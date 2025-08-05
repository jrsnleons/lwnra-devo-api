#!/bin/bash

# Facebook Permanent Page Access Token Generator
# This script helps you get a permanent (never-expiring) page access token

echo "🔑 Facebook Permanent Page Access Token Generator"
echo "==============================================="
echo

# Check if required tools are available
if ! command -v curl &> /dev/null; then
    echo "❌ curl is required but not installed."
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "⚠️  jq is recommended for pretty JSON output but not required."
    echo
fi

echo "📋 You'll need:"
echo "1. Your Facebook App ID: 307024831662935"
echo "2. Your Facebook App Secret: 042179ed1c7573fa33605f89d130e20e"
echo "3. A USER Access Token (not page token) from Graph API Explorer"
echo
echo "🔗 Get your USER token here:"
echo "   https://developers.facebook.com/tools/explorer/"
echo "   - Select your app: 'Website scrapper'"
echo "   - Make sure Token Type is 'User Access Token'"
echo "   - Required permissions: pages_read_engagement, pages_show_list"
echo
echo "⚠️  IMPORTANT: Make sure you select 'User Access Token' not 'Page Access Token'"
echo

# Get input from user
read -p "Enter your USER Access Token: " USER_TOKEN

if [ -z "$USER_TOKEN" ]; then
    echo "❌ User token is required"
    exit 1
fi

APP_ID="307024831662935"
APP_SECRET="042179ed1c7573fa33605f89d130e20e"

echo
echo "🔄 Step 1: Getting your managed pages..."

# Get pages managed by the user
PAGES_RESPONSE=$(curl -s "https://graph.facebook.com/v18.0/me/accounts?access_token=$USER_TOKEN")

# Check if jq is available for pretty printing
if command -v jq &> /dev/null; then
    echo "📄 Your pages:"
    echo "$PAGES_RESPONSE" | jq -r '.data[] | "ID: \(.id) | Name: \(.name) | Token expires: \(if .access_token then "Has token" else "No token" end)"'

    # Extract the Living Word NRA page token
    LWNRA_TOKEN=$(echo "$PAGES_RESPONSE" | jq -r '.data[] | select(.name == "Living Word NRA") | .access_token')
    LWNRA_ID=$(echo "$PAGES_RESPONSE" | jq -r '.data[] | select(.name == "Living Word NRA") | .id')
else
    echo "📄 Pages response:"
    echo "$PAGES_RESPONSE"
    echo
    echo "Please find the 'Living Word NRA' page in the response above and copy its access_token"
    read -p "Enter the Living Word NRA page access token: " LWNRA_TOKEN
    LWNRA_ID="164421594332429"
fi

if [ -z "$LWNRA_TOKEN" ] || [ "$LWNRA_TOKEN" = "null" ]; then
    echo "❌ Could not find Living Word NRA page token"
    echo "Make sure:"
    echo "1. You're an admin of the Living Word NRA page"
    echo "2. Your user token has pages_show_list permission"
    echo "3. You selected 'User Access Token' not 'Page Access Token'"
    exit 1
fi

echo
echo "✅ Found Living Word NRA page token"
echo "📄 Page ID: $LWNRA_ID"

echo
echo "🔄 Step 2: Verifying the page token..."

# Debug the page token to check its expiration
TOKEN_DEBUG=$(curl -s "https://graph.facebook.com/v18.0/debug_token?input_token=$LWNRA_TOKEN&access_token=$USER_TOKEN")

if command -v jq &> /dev/null; then
    EXPIRES_AT=$(echo "$TOKEN_DEBUG" | jq -r '.data.expires_at // "never"')
    TOKEN_TYPE=$(echo "$TOKEN_DEBUG" | jq -r '.data.type // "unknown"')

    echo "🔍 Token Info:"
    echo "   Type: $TOKEN_TYPE"
    echo "   Expires: $EXPIRES_AT"

    if [ "$EXPIRES_AT" = "0" ] || [ "$EXPIRES_AT" = "never" ]; then
        echo "🎉 SUCCESS! This token is already permanent (never expires)"
        echo
        echo "🔑 Your permanent Facebook page access token:"
        echo "$LWNRA_TOKEN"
        echo
        echo "💾 To use this token, set the environment variable:"
        echo "export FB_ACCESS_TOKEN=\"$LWNRA_TOKEN\""
        echo
        echo "🚀 For Railway deployment, add this environment variable:"
        echo "FB_ACCESS_TOKEN=$LWNRA_TOKEN"
    else
        echo "⚠️  This token still has an expiration date"
        echo "This usually means the page token is already as permanent as it can be"
        echo "for your current app setup."
        echo
        echo "🔑 Your Facebook page access token:"
        echo "$LWNRA_TOKEN"
        echo
        echo "💾 To use this token, set the environment variable:"
        echo "export FB_ACCESS_TOKEN=\"$LWNRA_TOKEN\""
        echo
        echo "🚀 For Railway deployment, add this environment variable:"
        echo "FB_ACCESS_TOKEN=$LWNRA_TOKEN"
    fi
else
    echo "Token debug response:"
    echo "$TOKEN_DEBUG"
fi

echo
echo "✅ Setup complete!"
echo
echo "📚 For more info about permanent tokens:"
echo "   https://developers.facebook.com/docs/pages/access-tokens#get-a-long-lived-page-access-token"
