# s3mounter
S3 mount to Local FS

# Setup for AWS-IAM

- Create User for IAM
- Generate Access Key Credentials

# Install AWS CLI

- sudo apt install -y curl unzip python-is-python3
- curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
- unzip awscli-bundle.zip
- sudo ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws
- sudo aws configure
- (input credentials)

# Build & Install

- task build install

# Direct Install

- go install github.com/nobonobo/s3mounter@latest

# Create Bucket and Root folder at S3-Console

- Create Bucket: "global-unique-bucket-name"
- Create Folder: "root/"

# Run

- mkdir dest
- sudo 3smounter -region ap-northeast-3 -bucket global-unique-bucket-name -root root dest
