awsctx
===========================

## Setup
Add your `.bashrc` or `.zshrc` and follow "OPTIONAL"
```bash
sync-aws-profile() {
    if [ -e "$AWSCTX/aws_profile" ]; then
        export AWS_PROFILE=$(cat $AWSCTX/aws_profile)
    fi
}
sync-aws-profile

# ---
# Enable tab completion

# zsh:
eval "$(awsctx completion zsh)"

# bash:
eval "$(awsctx completion bash)"
# ---

# ---
# OPTIONAL: Insert MFA Secret
# Uncomment me If you want to use CLI based 2FA
# MFA_TOKEN=$(oathtool -b --totp $AWS_MFA_SECRET)
# alias awsctx='awsctx --mfa=$(oathtool -b --totp $MFA_TOKEN)'
# ---
```

## Usage
```bash
awsctx [completion|help|list|login|show] ...

# List available accounts
awsctx list

# Show current AWS account
awsctx show

# Login
awsctx login myservice-production

# Login with MFA
awsctx login --mfa=123456 myservice-production

# Login and open AWS Management Console (-c, --console)
awsctx login -c myservice-production

# Login and get login URL (-l, --link)
awsctx login -l myservice-production

# Login and open AWS Management Console (with Firefox Multi-Account Containers) (-F, --firefox)
awsctx login -cF myservice-production
```
## Extra: Setup Firefox (Multi-Account Containers)
### Install addons:
Multi-Account Containers:
https://addons.mozilla.org/ja/firefox/addon/multi-account-containers/

Open container tabs from URL (Protocol Handler):
https://addons.mozilla.org/en-US/firefox/addon/open-url-in-container/

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
