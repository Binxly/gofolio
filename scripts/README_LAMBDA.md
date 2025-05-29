# AWS Lambda Deployment Guide

This guide explains how to deploy your Go + Templ portfolio website to AWS Lambda.

## Prerequisites

1. AWS CLI installed and configured with appropriate credentials
2. An AWS account with permissions to create Lambda functions and API Gateway
3. Go 1.21+ installed
4. Make installed

## Architecture

The application uses [algnhsa](https://github.com/akrylysov/algnhsa) to adapt the standard `net/http` handlers for AWS Lambda. This allows the same codebase to run both locally and on Lambda with minimal changes.

## Build Tags

- `main.go` - Used for local development (has build tag `!lambda`)
- `main_lambda.go` - Used for Lambda deployment (has build tag `lambda`)

## Building for Lambda

### For x86_64 architecture:
```bash
make build-lambda
```

### For ARM64 architecture (more cost-effective):
```bash
make build-lambda-arm64
```

This will:
1. Compile the Go application with Lambda build tags
2. Create a `bootstrap` executable (required by Lambda custom runtime)
3. Package everything into `function.zip` including static files, assets, and content

## Deployment Options

### Option 1: Lambda Function URL (Simplest)

1. Create the Lambda function:
```bash
AWS_ACCOUNT_ID=your-account-id make create-function
```

2. In AWS Console:
   - Go to your Lambda function
   - Navigate to "Configuration" → "Function URL"
   - Click "Create function URL"
   - Choose "NONE" for auth type (or configure as needed)
   - Save the function URL

### Option 2: API Gateway HTTP API

1. Create Lambda function as above
2. Create HTTP API in API Gateway:
   - Choose "HTTP API"
   - Add Lambda integration
   - Configure routes: `$default` route pointing to your Lambda

### Option 3: API Gateway REST API

1. Create Lambda function as above
2. Create REST API:
   - Enable "Lambda Proxy Integration"
   - Create `ANY` method on `/`
   - Create `{proxy+}` resource with `ANY` method

## Static Assets

**Important**: For production, static assets (CSS, JS, images) should be served from S3 + CloudFront rather than Lambda. Lambda is not optimized for serving static files.

### Setting up S3 + CloudFront:

1. Create S3 bucket for static assets
2. Upload `static/` and `assets/` directories
3. Create CloudFront distribution pointing to S3
4. Update your templates to use CloudFront URLs

## Environment Variables

Set these in Lambda configuration:
- `PORT` is ignored (Lambda assigns its own)
- Add any other environment variables your app needs

## Updating Deployed Function

After making changes:
```bash
make deploy
```

## Cost Optimization

1. Use ARM64 architecture (Graviton2) - up to 34% better price/performance
2. Set appropriate memory size (start with 512MB)
3. Enable Lambda SnapStart for faster cold starts
4. Use CloudFront for static assets
5. Consider using Lambda@Edge for global distribution

## Monitoring

- Enable CloudWatch Logs
- Set up CloudWatch Alarms for errors
- Use AWS X-Ray for tracing
- Monitor cold start frequency

## Limitations

1. Lambda has a 15-minute timeout limit
2. Response payload limit is 6MB
3. Request payload limit is 10MB
4. Temporary storage (`/tmp`) is limited to 512MB

## Troubleshooting

### Function returns 502 Bad Gateway
- Check CloudWatch Logs for errors
- Ensure all file paths are relative
- Verify all dependencies are included in the zip

### Static files not loading
- Static files should be served from S3/CloudFront in production
- For testing, ensure they're included in the zip file

### Cold start issues
- Consider using Lambda SnapStart
- Increase memory allocation
- Use provisioned concurrency for consistent performance 