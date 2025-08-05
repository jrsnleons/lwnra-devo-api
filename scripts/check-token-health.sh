#!/bin/bash

# Token Monitoring and Management Script
# Checks token health and provides renewal instructions

echo "🔍 Facebook Token Health Monitor"
echo "==============================="
echo

# Configuration
CURRENT_TOKEN="EAAEXPMoHS1cBPPULXqZB9Qmnrm7Uqk7gq3nIWe1ROPTODyBIqroV5b8CZCp8qVsTjnaDbdew2uFx69huv2H9T8ZBsk5CZBefq7d2PvdZCSuAol040217WSr2uKNI29cM6CezwfuK5e5UgAd2hiiO8DyeIqaqxnfwZC1iQwF6aZCZBDFbJGueQMlpGNiqhcZCLcn4mmql8pQwAbJ0jlatAkqY5nwneucZCFuYYXBJdtaQZDZD"
PAGE_ID="164421594332429"

echo "📊 Checking current token status..."

# Get token info
TOKEN_INFO=$(curl -s "https://graph.facebook.com/v18.0/debug_token?input_token=$CURRENT_TOKEN&access_token=$CURRENT_TOKEN")

if command -v jq &> /dev/null; then
    echo "🔍 Token Details:"

    IS_VALID=$(echo "$TOKEN_INFO" | jq -r '.data.is_valid // false')
    TOKEN_TYPE=$(echo "$TOKEN_INFO" | jq -r '.data.type // "unknown"')
    EXPIRES_AT=$(echo "$TOKEN_INFO" | jq -r '.data.expires_at // 0')
    APP_ID=$(echo "$TOKEN_INFO" | jq -r '.data.app_id // "unknown"')

    echo "   ✅ Valid: $IS_VALID"
    echo "   📄 Type: $TOKEN_TYPE"
    echo "   🆔 App ID: $APP_ID"

    if [ "$EXPIRES_AT" != "0" ] && [ "$EXPIRES_AT" != "null" ]; then
        EXPIRES_DATE=$(date -r "$EXPIRES_AT" 2>/dev/null || echo "Invalid timestamp")
        CURRENT_TIME=$(date +%s)
        DAYS_LEFT=$(( ($EXPIRES_AT - $CURRENT_TIME) / 86400 ))

        echo "   ⏰ Expires: $EXPIRES_DATE"
        echo "   📅 Days left: $DAYS_LEFT"

        if [ "$DAYS_LEFT" -lt 7 ]; then
            echo "   ⚠️  WARNING: Token expires soon!"
        elif [ "$DAYS_LEFT" -lt 30 ]; then
            echo "   🟡 Token expires in less than 30 days"
        else
            echo "   ✅ Token has good expiration time"
        fi
    else
        echo "   ♾️  Token: Never expires (permanent)"
    fi

    echo
    echo "🧪 Testing API access..."

    # Test API call
    API_TEST=$(curl -s "https://graph.facebook.com/v18.0/$PAGE_ID?fields=id,name&access_token=$CURRENT_TOKEN")
    API_SUCCESS=$(echo "$API_TEST" | jq -r '.name // "ERROR"')

    if [ "$API_SUCCESS" != "ERROR" ]; then
        echo "   ✅ API Access: Working"
        echo "   📄 Page: $API_SUCCESS"
    else
        echo "   ❌ API Access: Failed"
        echo "   Error: $(echo "$API_TEST" | jq -r '.error.message // "Unknown error"')"
    fi

else
    echo "Raw token info:"
    echo "$TOKEN_INFO"
fi

echo
echo "🔄 Token Renewal Instructions:"
echo "1. Go to Facebook Graph API Explorer: https://developers.facebook.com/tools/explorer/"
echo "2. Select your app: 'Website scrapper'"
echo "3. Make sure to select 'User Access Token' (not Page)"
echo "4. Add permissions: pages_read_engagement, pages_show_list"
echo "5. Generate token and run the permanent token script"
echo
echo "📝 Quick renewal command:"
echo "   ./scripts/get-permanent-token.sh"
echo
echo "🚀 For production, consider setting up automated token refresh"
echo "   or use Facebook System Users for truly permanent tokens"
