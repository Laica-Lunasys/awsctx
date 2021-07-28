#!/bin/bash
set -e
cd "$(dirname "$0")"

_TARGET="$1"
_MFA_TOKEN="$2"

_login() {
    _account="$1"
    _mfa="$2"
    echo ":: Logging in to: $_account"
    saml2aws login -a $_account \
        --skip-prompt \
        --mfa-token="$_mfa"

    echo ":: Update AWS_PROFILE: $_account"
    touch $PWD/aws_profile
    echo "$_account" > $PWD/aws_profile
}

if [ "$_TARGET" == "all" ]; then
    for e in $(cat ~/.saml2aws | grep "\[*\]" | cat); do
        _a=$(echo $e | tr -d '[]')
        _login $_a $_MFA_TOKEN
    done
elif [ "$_TARGET" != "" ] && [ "$_MFA_TOKEN" != "" ]; then
    _login $_TARGET $_MFA_TOKEN
else
    echo "Usage: awsctx <target/all> <mfa_token>"
    exit 1
fi
