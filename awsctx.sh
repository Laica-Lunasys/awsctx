#!/bin/bash
set -e
cd "$(dirname "$0")"

_TARGET="$1"
_MFA_TOKEN="$2"
_OPTIONS="$3"

_login() {
    _account="$1"
    _mfa="$2"
    _options="$3"
    echo ":: Logging in to: $_account"
    saml2aws login -a $_account \
        --skip-prompt \
        --mfa-token="$_mfa"

    echo ":: Update AWS_PROFILE: $_account"
    touch $PWD/aws_profile
    echo "$_account" > $PWD/aws_profile

    if [ "$_options" != "" ]; then
        saml2aws $_options $(shift 3 && echo $@) -a $_account \
            --skip-prompt \
            --mfa-token="$_mfa"
    fi
}

if [ "$_TARGET" == "all" ]; then
    for e in $(cat ~/.saml2aws | grep "\[*\]" | cat); do
        _a=$(echo $e | tr -d '[]')
        _login $_a $_MFA_TOKEN $_OPTIONS $(shift 3 && echo $@)
    done
elif [ "$_TARGET" != "" ] && [ "$_MFA_TOKEN" != "" ]; then
    _login $_TARGET $_MFA_TOKEN $_OPTIONS $(shift 3 && echo $@)
else
    echo "Usage: awsctx <target/all> [MFA_TOKEN] [login|list-roles|console]"
    exit 1
fi
