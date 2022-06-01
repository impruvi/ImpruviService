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
    this.createS3Bucket('impruvi-drills');
    this.createS3Bucket('impruvi-submissions');
    this.createS3Bucket('impruvi-feedback');
  }

  createIAMRole = (domain: string) => {
    return new iam.Role(this, `${domain}-BentoServiceRole`, {
      roleName: `${domain}-BentoServiceRole`,
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/CloudWatchLogsFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSQSFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AWSLambda_FullAccess'}
      ]
    });
  };

  createDynamoTables = () => {
    const usersTable = new dynamodb.Table(this, `${this.domain}-users`, {
      partitionKey: { name: 'userId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-users`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    usersTable.addGlobalSecondaryIndex({
      indexName: 'invitation-code-index',
      partitionKey: {name: 'invitationCode', type: dynamodb.AttributeType.STRING},
    });

    new dynamodb.Table(this, `${this.domain}-coaches`, {
      partitionKey: { name: 'coachId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-coaches`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    new dynamodb.Table(this, `${this.domain}-sessions`, {
      partitionKey: { name: 'userId', type: dynamodb.AttributeType.STRING},
      sortKey: { name: 'sessionNumber', type: dynamodb.AttributeType.NUMBER},
      tableName: `${this.domain}-sessions`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    new dynamodb.Table(this, `${this.domain}-drills`, {
      partitionKey: { name: 'drillId', type: dynamodb.AttributeType.STRING},
      tableName: `${this.domain}-drills`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
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

  createApiResources = (iamRole: any) => {
    const apiHandlerLambda = new lambda.Function(this, `${this.domain}-impruvi-service-api-handler`, {
      functionName: `${this.domain}-impruvi-service-api-handler`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.seconds(8),
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
        ['/validate-invitation-code', [HttpMethod.POST]],
        ['/get-sessions', [HttpMethod.POST]],
        ['/get-video-upload-url', [HttpMethod.POST]],
        ['/create-submission', [HttpMethod.POST]],
        ['/create-feedback', [HttpMethod.POST]],
        ['/get-all-users', [HttpMethod.POST]],
        ['/get-all-drills', [HttpMethod.POST]],
        ['/update-session', [HttpMethod.POST]],
      ])
    });
  };
}
