#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { ImpruviServiceStack } from '../lib/impruvi_service-stack';

const app = new cdk.App();
new ImpruviServiceStack(app, 'ImpruviServiceStack');
