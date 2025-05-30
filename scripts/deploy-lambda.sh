#!/bin/bash

# Lambda Deployment Script
set -e

# Configuration
FUNCTION_NAME="${LAMBDA_FUNCTION_NAME:-function}"
AWS_REGION="${AWS_REGION:-us-east-1}"
ARCHITECTURE="${LAMBDA_ARCH:-arm64}" # arm64 is more cost-effective

echo "🚀 Building Lambda function..."

# Generate templ files
echo "📝 Generating templ files..."
go tool templ generate

# Build CSS
echo "🎨 Building CSS..."
tailwindcss -i ./static/assets/css/input.css -o ./static/assets/css/output.css

# Build based on architecture
if [ "$ARCHITECTURE" = "arm64" ]; then
    echo "🔨 Building for ARM64..."
    GOOS=linux GOARCH=arm64 go build -tags lambda,lambda.norpc -o bootstrap .
else
    echo "🔨 Building for x86_64..."
    GOOS=linux GOARCH=amd64 go build -tags lambda,lambda.norpc -o bootstrap .
fi

# Create deployment package
echo "📦 Creating deployment package..."
zip -r function.zip bootstrap

# Add static assets to the package
echo "📁 Adding static assets..."
zip -r function.zip static/ assets/ blog/posts/

# Get file size
FILE_SIZE=$(stat -f%z function.zip 2>/dev/null || stat -c%s function.zip)
echo "📊 Package size: $(($FILE_SIZE / 1024 / 1024))MB"

# Deploy or update
echo "🔄 Updating function code..."
aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://function.zip \
    --region $AWS_REGION

if [ $? -eq 0 ]; then
    echo "✅ Deployment successful!"
    echo "🔗 Your site: https://binx.page"
else
    echo "❌ Deployment failed!"
    exit 1
fi

# Cleanup
echo "🧹 Cleaning up..."
rm -f bootstrap function.zip

echo "🎉 All done!" 