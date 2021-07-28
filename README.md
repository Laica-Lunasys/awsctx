awsctx
===========================

## Setup
Add your `.bashrc` or `.zshrc` and follow "REPLACE ME" and "OPTIONAL"
```bash
AWSCTX="$HOME/awsctx" # <-- REPLACE ME: path to awsctx
sync-aws-profile() {
    if [ -e "$AWSCTX/aws_profile" ]; then
        export AWS_PROFILE=$(cat $AWSCTX/aws_profile)
    fi
}
awsctx() {
    # ---
    # OPTIONAL: Insert MFA Secret
    # Uncomment me If you want to use CLI based 2FA
    # MFA_TOKEN=$(oathtool -b --totp $AWS_MFA_SECRET)
    # ---

    if [ "$MFA_TOKEN" != "" ]; then
        bash $AWSCTX/awsctx.sh $@ $MFA_TOKEN
    else
        bash $AWSCTX/awsctx.sh $@
    fi
    sync-aws-profile
}
sync-aws-profile
```

## Usage
```bash
awsctx <AWS_ACCOUNT> [MFA_TOKEN]
# If you defined MFA_TOKEN: `awsctx myservice-production`
# or else: `awsctx myservice-production 123456`
```

## Extra: Setup oathtool (MFA_TOKEN)

### twilio-authy

1. Export TOTP tokens
Get authy-export from: https://github.com/alexzorin/authy
```bash
authy-export
# Enter TOTP backup password
```

2. Copy secret
```
otpauth://totp/~~~:~~~?digits=~&secret=<COPY_HERE>
```

Save to env
```bash
AWS_MFA_SECRET=<SECRET>
```

### Get token from oathtool

1. Install oath-toolkit
```bash
# Ubuntu/Debian
sudo apt install oathtool

# Fedora/RHEL
sudo dnf install oathtool

# macOS
brew install oath-toolkit
```

2. Get TOKEN from SECRET
```bash
oathtool -b --totp $AWS_MFA_SECRET
```
