import * as cdk from '@aws-cdk/core';
import * as iam from '@aws-cdk/aws-iam';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as events from '@aws-cdk/aws-events';
import * as eventsTargets from '@aws-cdk/aws-events-targets';
import * as lambda from '@aws-cdk/aws-lambda';
import {HttpMethod, SingleLambdaBackedRestApi} from '@climatehub/cdk-constructs';
import * as s3 from '@aws-cdk/aws-s3';
import * as cloudfront from '@aws-cdk/aws-cloudfront';
import * as cloudfrontOrigins from '@aws-cdk/aws-cloudfront-origins';
import * as certificateManager from '@aws-cdk/aws-certificatemanager';
import * as stepFunction from '@aws-cdk/aws-stepfunctions';
import * as stepFunctionTasks from '@aws-cdk/aws-stepfunctions-tasks';
import * as mediaconvert from '@aws-cdk/aws-mediaconvert';


const path = require('path');

const getStripeSecretKey = (domain: string) => {
  switch (domain) {
    case 'beta':
      return 'sk_test_51LIhrlKA3EgJIYsfR79B9PLo9RVRXr66oAL70oOO8XUZARIk2QTCkM3vKXdm7Bp4oo9T8aRrFEj6kvroWsndlM7F00c5h6D8YY';
    default:
      return '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'
  }
}

const getWebhookSigningSecret = (domain: string) => {
  switch (domain) {
    case 'beta':
      return 'whsec_1GbPJpu2ibnLMJixHPcBuBSvOWiEO8Qm';
    default:
      return 'whsec_y2tsLGSix8tE8MRDEfMHCL6ncy5CMSnW'
  }
}

export interface ImpruviServiceStackProps extends cdk.StackProps {
  readonly domain: string;
  readonly env: cdk.Environment
}

export class ImpruviServiceStack extends cdk.Stack {
  private readonly domain: string;

  constructor(scope: cdk.Construct, id: string, props: ImpruviServiceStackProps) {
    super(scope, id, props);
    this.domain = props.domain;

    const iamRole = this.createIAMRole();
    this.createDynamoTables();
    this.createApiResources(iamRole);
    this.createFixedReminderLambda(iamRole);
    this.createMediaBucket();
    this.createCloudfrontDistribution();
    this.createReminderStepFunction(iamRole);
    this.createMediaConvertQueue();
    this.createMediaConvertHandlerLambda(iamRole);
  }

  createMediaConvertQueue = () => {
    new mediaconvert.CfnQueue(this, `${this.domain}-impruvi-service-queue`, {
      name: `${this.domain}-impruvi-service-queue`,
      description: 'queue for transcoding videos',
      pricingPlan: 'ON_DEMAND',
    })
  }

  createMediaConvertHandlerLambda = (iamRole: any) => {
    const eventHandler = new lambda.Function(this, `${this.domain}-impruvi-service-mediaconvert-event-handler`, {
      functionName: `${this.domain}-impruvi-service-mediaconvert-event-handler`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.minutes(5),
      environment: {
        DOMAIN: this.domain,
        STRIPE_SECRET_KEY: getStripeSecretKey(this.domain),
        WEB_HOOK_SIGNING_SECRET: getWebhookSigningSecret(this.domain)
      },
      tracing: lambda.Tracing.ACTIVE
    });

    new events.Rule(this, `${this.domain}-impruvi-service-mediaconvert-event`, {
      ruleName: `${this.domain}-impruvi-service-mediaconvert-event`,
      eventPattern: {
        source: ["aws.mediaconvert"],
        detailType: ["MediaConvert Job State Change"],
        detail: {
          status: ["ERROR", "COMPLETE"]
        }
      },
      targets: [
        new eventsTargets.LambdaFunction(eventHandler, {})
      ],
    });
  }

  createIAMRole = () => {
    return new iam.Role(this, `${this.domain}-ImpruviServiceRole`, {
      roleName: `${this.domain}-ImpruviServiceRole`,
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/CloudWatchLogsFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSQSFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AWSLambda_FullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonS3FullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSNSFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AmazonSESFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AWSStepFunctionsFullAccess'},
        {managedPolicyArn: 'arn:aws:iam::aws:policy/AWSElementalMediaConvertFullAccess'}
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
    playersTable.addGlobalSecondaryIndex({
      indexName: 'email-index',
      partitionKey: {name: 'email', type: dynamodb.AttributeType.STRING},
    });

    new dynamodb.Table(this, `${this.domain}-password-reset-codes`, {
      partitionKey: { name: 'email', type: dynamodb.AttributeType.STRING},
      sortKey: { name: 'creationDateEpochMillis', type: dynamodb.AttributeType.NUMBER},
      tableName: `${this.domain}-password-reset-codes`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
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

  createMediaBucket = () => {
    const bucket = new s3.Bucket(this, `${this.domain}-impruvi-media-bucket`, {
      bucketName: `${this.domain}-impruvi-media-bucket`,
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

    const bucketPolicy = new s3.BucketPolicy(this, `${this.domain}-impruvi-media-policy`, {
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

  createReminderStepFunction = (iamRole: any) => {
    const reminderNotificationWait = new stepFunction.Wait(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-wait`, {
      time: stepFunction.WaitTime.secondsPath('$.waitSeconds'),
    });

    const notificationSender = new lambda.Function(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-sender`, {
      functionName: `${this.domain}-impruvi-service-dynamic-reminder-notification-sender`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.minutes(5),
      environment: {
        DOMAIN: this.domain,
        STRIPE_SECRET_KEY: getStripeSecretKey(this.domain),
        WEB_HOOK_SIGNING_SECRET: getWebhookSigningSecret(this.domain)
      },
      tracing: lambda.Tracing.ACTIVE
    });

    const reminderNotificationSenderTask = new stepFunctionTasks.LambdaInvoke(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-sender-task`, {
      lambdaFunction: notificationSender,
      outputPath: '$.Payload',
    });

    const reminderNotificationPassthrough = new stepFunction.Pass(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-pass`)

    const reminderNotificationChoice = new stepFunction.Choice(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-choice`)
        .when(stepFunction.Condition.not(stepFunction.Condition.booleanEquals('$.completed', true)), reminderNotificationWait)
        .otherwise(reminderNotificationPassthrough);

    const definition = reminderNotificationWait
        .next(reminderNotificationSenderTask)
        .next(reminderNotificationChoice);

    new stepFunction.StateMachine(this, `${this.domain}-impruvi-service-dynamic-reminder-notification-state-machine`, {
      definition,
      stateMachineName: `${this.domain}-impruvi-service-dynamic-reminder-notification-state-machine`,
      timeout: cdk.Duration.days(1),
    });
  }

  createFixedReminderLambda = (iamRole: any) => {
    const notificationSender = new lambda.Function(this, `${this.domain}-impruvi-service-fixed-reminder-notification-sender`, {
      functionName: `${this.domain}-impruvi-service-fixed-reminder-notification-sender`,
      runtime: lambda.Runtime.GO_1_X,
      handler: 'ImpruviService',
      role: iamRole,
      code: lambda.Code.fromAsset(path.join(__dirname, '/build')),
      memorySize: 2048,
      timeout:  cdk.Duration.minutes(15),
      environment: {
        DOMAIN: this.domain,
        STRIPE_SECRET_KEY: getStripeSecretKey(this.domain),
        WEB_HOOK_SIGNING_SECRET: getWebhookSigningSecret(this.domain)
      },
      tracing: lambda.Tracing.ACTIVE
    });
    new events.Rule(this, `${this.domain}-impruvi-service-fixed-notification-sender-rule`, {
      ruleName: `${this.domain}-impruvi-service-fixed-notification-sender-rule`,
      schedule: events.Schedule.cron({
        weekDay: '2,5',
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
        DOMAIN: this.domain,
        STRIPE_SECRET_KEY: getStripeSecretKey(this.domain),
        WEB_HOOK_SIGNING_SECRET: getWebhookSigningSecret(this.domain)
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

        ['/subscription-plan/get', [HttpMethod.POST]],

        ['/player/sign-up/initiate', [HttpMethod.POST]],
        ['/player/sign-up/complete', [HttpMethod.POST]],
        ['/player/sign-in', [HttpMethod.POST]],
        ['/player/payment-methods/get', [HttpMethod.POST]],
        ['/player/subscription/re-activate', [HttpMethod.POST]],
        ['/player/subscription/get', [HttpMethod.POST]],
        ['/player/subscription/create', [HttpMethod.POST]],
        ['/player/subscription/cancel', [HttpMethod.POST]],
        ['/player/password-reset/initiate', [HttpMethod.POST]],
        ['/player/password-reset/validate-code', [HttpMethod.POST]],
        ['/player/password-reset/complete', [HttpMethod.POST]],
        ['/player/update', [HttpMethod.POST]],
        ['/player/get', [HttpMethod.POST]],
        ['/player/inbox/get', [HttpMethod.POST]],

        ['/coaches/list', [HttpMethod.POST]],
        ['/coach/update', [HttpMethod.POST]],
        ['/coach/get', [HttpMethod.POST]],
        ['/coach/players-and-subscriptions/get', [HttpMethod.POST]],

        ['/sessions/player/get', [HttpMethod.POST]],
        ['/sessions/coach/get', [HttpMethod.POST]],
        ['/sessions/get', [HttpMethod.POST]],
        ['/sessions/delete', [HttpMethod.POST]],
        ['/sessions/create', [HttpMethod.POST]],
        ['/sessions/update', [HttpMethod.POST]],
        ['/sessions/submission/create', [HttpMethod.POST]],
        ['/sessions/feedback/create', [HttpMethod.POST]],
        ['/sessions/feedback/view', [HttpMethod.POST]],

        ['/drills/get', [HttpMethod.POST]],
        ['/drills/create', [HttpMethod.POST]],
        ['/drills/update', [HttpMethod.POST]],
        ['/drills/delete', [HttpMethod.POST]],
        ['/drills/coach/get', [HttpMethod.POST]],
        ['/drills/player/get', [HttpMethod.POST]],

        ['/media-upload-url/generate', [HttpMethod.POST]],

        ['/stripe-event', [HttpMethod.POST]],

        ['/app-version/is-compatible', [HttpMethod.POST]],
      ])
    });
  };

  createCloudfrontDistribution = () => {
    const bucket = new s3.Bucket(this, `${this.domain}-impruvi-web-static-assets`, {
      bucketName: `${this.domain}-impruvi-web-static-assets`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      cors: [
        {
          allowedOrigins: ['*'],
          allowedMethods: [s3.HttpMethods.HEAD, s3.HttpMethods.GET],
          allowedHeaders: ['*']
        }
      ]
    });

    const originAccessIdentity = new cloudfront.OriginAccessIdentity(this, `${this.domain}-impruvi-web-origin-access-identity`, {});
    const bucketPolicy = new s3.BucketPolicy(this, `${this.domain}-impruvi-web-assets-policy`, {
      bucket: bucket
    });

    bucketPolicy.document.addStatements(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: ['s3:GetObject'],
      principals: [
        new iam.AnyPrincipal(),
        new iam.CanonicalUserPrincipal(originAccessIdentity.cloudFrontOriginAccessIdentityS3CanonicalUserId),
      ],
      resources: [
        bucket.bucketArn + '/*'
      ],
    }));

    const domainNames = this.domain === 'prod' ? ['impruviapp.com'] : undefined;
    const certificate = this.domain === 'prod'
        ? certificateManager.Certificate.fromCertificateArn(this, "sslCertificate", "arn:aws:acm:us-east-1:522042996447:certificate/8e8a4051-4063-4faa-9b47-db7999e9ad35")
        : undefined;
    new cloudfront.Distribution(this, `${this.domain}-impruvi-web-distribution`, {
      defaultBehavior: {
        origin: new cloudfrontOrigins.S3Origin(bucket, {
          originAccessIdentity: originAccessIdentity
        }),
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS
      },
      defaultRootObject: 'index.html',
      errorResponses: [
        {
          httpStatus: 403,
          responsePagePath: '/index.html'
        }
      ],
      domainNames: domainNames,
      certificate: certificate
    });
  }
}
