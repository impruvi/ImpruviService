#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import {ImpruviServiceStack} from "../lib/impruvi_service-stack";

const app = new cdk.App();
const domain = process.env.DOMAIN;
if (!domain) {
    console.error("a domain must be provided");
    process.exit(1);
}
const env = { account: '730511296908', region: 'us-west-2' };

new ImpruviServiceStack(app, `${domain}-ImpruviServiceStack`, {
    domain: domain,
    env: env
});
