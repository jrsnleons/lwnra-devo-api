#!/bin/bash

# Facebook System User Token Generator (For Production Apps)
# This creates a permanent token using Facebook System Users

echo "🤖 Facebook System User Token Generator"
echo "======================================"
echo

echo "📋 This method creates a truly permanent token using Facebook System Users"
echo "🏢 Recommended for production applications that need continuous access"
echo

echo "⚠️  PREREQUISITES:"
echo "1. Your Facebook app must be in 'Live' mode (not Development)"
echo "2. You must be an admin of the Facebook app"
echo "3. You need your App ID and App Secret"
echo

echo "🔗 Steps to create a System User:"
echo "1. Go to: https://business.facebook.com/settings/system-users"
echo "2. Click 'Add' to create a new System User"
echo "3. Give it a name like 'LWNRA Devotional API'"
echo "4. Assign it to your Facebook app"
echo "5. Generate an access token for the System User"
echo

echo "📖 After creating the System User:"
echo "1. The token generated will be permanent (no expiration)"
echo "2. You can assign specific permissions"
echo "3. It's designed for server-to-server communication"
echo

echo "🚀 Alternative: Use App Access Token for basic operations"

read -p "Do you want to try generating an App Access Token instead? (y/n): " USE_APP_TOKEN

if [ "$USE_APP_TOKEN" = "y" ] || [ "$USE_APP_TOKEN" = "Y" ]; then
    APP_ID="307024831662935"
    APP_SECRET="042179ed1c7573fa33605f89d130e20e"

    echo
    echo "🔄 Generating App Access Token..."

    APP_TOKEN="${APP_ID}|${APP_SECRET}"

    echo "✅ App Access Token generated:"
    echo "$APP_TOKEN"
    echo
    echo "⚠️  NOTE: App Access Tokens have limited permissions"
    echo "They work for some API calls but may not work for all page operations"
    echo
    echo "🧪 Testing the App Access Token..."

    # Test the app token
    curl -s "https://graph.facebook.com/v18.0/164421594332429?access_token=$APP_TOKEN" | head -200

    echo
    echo "💾 To use this token:"
    echo "export FB_ACCESS_TOKEN=\"$APP_TOKEN\""
else
    echo
    echo "📖 For System User setup, visit:"
    echo "   https://business.facebook.com/settings/system-users"
    echo
    echo "📚 Documentation:"
    echo "   https://developers.facebook.com/docs/marketing-api/system-users"
fi
