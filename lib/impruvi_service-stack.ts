import * as cdk from '@aws-cdk/core';
import * as iam from '@aws-cdk/aws-iam';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as events from '@aws-cdk/aws-events';
import * as eventsTargets from '@aws-cdk/aws-events-targets';
import * as lambda from '@aws-cdk/aws-lambda';
import {HttpMethod, SingleLambdaBackedRestApi} from '@climatehub/cdk-constructs';
import * as s3 from '@aws-cdk/aws-s3';


const path = require('path');

export interface ImpruviServiceStackProps extends cdk.StackProps {
  readonly domain: string;
  readonly env: cdk.Environment
}

export class ImpruviServiceStack extends cdk.Stack {
  private readonly domain: string;

  constructor(scope: cdk.Construct, id: string, props: ImpruviServiceStackProps) {
    super(scope, id, props);
    this.domain = props.domain;

    const iamRole = this.createIAMRole(this.domain);
    this.createDynamoTables();
    this.createApiResources(iamRole);
    this.createAsyncLambdaResources(iamRole);
    this.createS3Bucket('impruvi-media');
  }

  createIAMRole = (domain: string) => {
    return new iam.Role(this, `${domain}-ImpruviServiceRole`, {
      roleName: `${domain}-ImpruviServiceRole`,
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/CloudWatchLogsFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSQSFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AWSLambda_FullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonS3FullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSNSFullAccess'}
      ]
    });
  };

  createDynamoTables = () => {
    const playersTable = new dynamodb.Table(this, `${this.domain}-players`, {
      partitionKey: { name: 'playerId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-players`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    playersTable.addGlobalSecondaryIndex({
      indexName: 'coachId-index',
      partitionKey: {name: 'coachId', type: dynamodb.AttributeType.STRING},
    });

    new dynamodb.Table(this, `${this.domain}-coaches`, {
      partitionKey: { name: 'coachId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-coaches`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    new dynamodb.Table(this, `${this.domain}-invitation-codes`, {
      partitionKey: { name: 'invitationCode', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-invitation-codes`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    new dynamodb.Table(this, `${this.domain}-sessions`, {
      partitionKey: { name: 'playerId', type: dynamodb.AttributeType.STRING},
      sortKey: { name: 'sessionNumber', type: dynamodb.AttributeType.NUMBER},
      tableName: `${this.domain}-sessions`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const drillsTable = new dynamodb.Table(this, `${this.domain}-drills`, {
      partitionKey: { name: 'drillId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-drills`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    drillsTable.addGlobalSecondaryIndex({
      indexName: 'coachId-index',
      partitionKey: {name: 'coachId', type: dynamodb.AttributeType.STRING},
    });
  };

  createS3Bucket = (bucketName: string) => {
    const bucket = new s3.Bucket(this, `${this.domain}-${bucketName}-bucket`, {
      bucketName: `${this.domain}-${bucketName}-bucket`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      cors: [
        {
          allowedHeaders: [
            "*"
          ],
          allowedMethods: [
            s3.HttpMethods.PUT,
            s3.HttpMethods.POST,
            s3.HttpMethods.GET,
            s3.HttpMethods.DELETE
          ],
          allowedOrigins: [
            "*"
          ],
          exposedHeaders: [
            "x-amz-server-side-encryption",
            "x-amz-request-id",
            "x-amz-id-2"
          ],
          maxAge: 3000
        }
      ]
    });

    const bucketPolicy = new s3.BucketPolicy(this, `${this.domain}-${bucketName}-policy`, {
      bucket: bucket
    });

    bucketPolicy.document.addStatements(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: ['s3:GetObject'],
      principals: [
        new iam.AnyPrincipal()
      ],
      resources: [
        bucket.bucketArn + '/*'
      ],
    }));
  };

  createAsyncLambdaResources = (iamRole: any) => {
    const notificationSender = new lambda.Function(this, `${this.domain}-impruvi-service-notification-sender`, {
      functionName: `${this.domain}-impruvi-service-notification-sender`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.minutes(15),
      environment: {
        domain: this.domain
      },
      tracing: lambda.Tracing.ACTIVE
    });
    new events.Rule(this, `${this.domain}-impruvi-service-notification-sender-rule`, {
      ruleName: `${this.domain}-impruvi-service-notification-sender-rule`,
      schedule: events.Schedule.cron({
        hour: '8',
        minute: '0',
      }),
      targets: [
        new eventsTargets.LambdaFunction(notificationSender, {
          event: events.RuleTargetInput.fromObject({
            body: "SEND_NOTIFICATIONS_EVENT"
          })
        })
      ],
    });
  }

  createApiResources = (iamRole: any) => {
    const apiHandlerLambda = new lambda.Function(this, `${this.domain}-impruvi-service-api-handler`, {
      functionName: `${this.domain}-impruvi-service-api-handler`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.seconds(10),
      environment: {
        domain: this.domain
      },
      tracing: lambda.Tracing.ACTIVE
    });

    new events.Rule(this, `${this.domain}-impruvi-service-api-handler-warm-up-rule`, {
      ruleName: `${this.domain}-impruvi-service-api-handler-warm-up-rule`,
      schedule: events.Schedule.rate(cdk.Duration.minutes(1)),
      targets: [
        new eventsTargets.LambdaFunction(apiHandlerLambda, {
          event: events.RuleTargetInput.fromObject({
            body: "WARM_UP_EVENT"
          })
        })
      ],
    });

    new SingleLambdaBackedRestApi(this, `${this.domain}-impruvi-service-api`, {
      restApiName: `${this.domain}-impruvi-service-api`,
      handler: apiHandlerLambda,
      corsEnabled: true,
      loggingOptions: {
        accessLoggingEnabled: true,
        dataTraceEnabled: true,
        metricsEnabled: true,
        tracingEnabled: true,
      },
      resources: new Map<string, HttpMethod[]>([
        ['/invitation-code/validate', [HttpMethod.POST]],

        ['/player/update', [HttpMethod.POST]],
        ['/player/get', [HttpMethod.POST]],
        ['/player/inbox/get', [HttpMethod.POST]],

        ['/coach/update', [HttpMethod.POST]],
        ['/coach/get', [HttpMethod.POST]],

        ['/sessions/player/get', [HttpMethod.POST]],
        ['/sessions/coach/get', [HttpMethod.POST]],
        ['/sessions/delete', [HttpMethod.POST]],
        ['/sessions/create', [HttpMethod.POST]],
        ['/sessions/update', [HttpMethod.POST]],
        ['/sessions/submission/create', [HttpMethod.POST]],
        ['/sessions/feedback/create', [HttpMethod.POST]],
        ['/sessions/feedback/view', [HttpMethod.POST]],

        ['/drills/create', [HttpMethod.POST]],
        ['/drills/update', [HttpMethod.POST]],
        ['/drills/delete', [HttpMethod.POST]],
        ['/drills/coach/get', [HttpMethod.POST]],
        ['/drills/player/get', [HttpMethod.POST]],

        ['/media-upload-url/generate', [HttpMethod.POST]],
      ])
    });
  };
}
