*** Settings ***
Resource         resources/common.resource
Suite Setup      Create Live Sessions
Test Setup       Prime Live Auction

*** Test Cases ***
API Root Exposes Service Metadata
    ${payload}=    Get JSON Response    api    /
    Dictionary Should Contain Item    ${payload}    service    those-forgotten-api
    Dictionary Should Contain Item    ${payload}    status    ok
    Dictionary Should Contain Key    ${payload}    swagger

API Health And Readiness Are Reachable
    ${health}=    Get JSON Response    api    /healthz
    Dictionary Should Contain Item    ${health}    status    ok
    ${ready}=    Get JSON Response    api    /readyz
    Dictionary Should Contain Item    ${ready}    status    ready

API Swagger Endpoints Are Reachable
    ${swagger_html}=    Get Text Response    api    /swagger
    Response Should Contain    ${swagger_html}    SwaggerUIBundle    /swagger.yaml
    ${swagger_yaml}=    Get Text Response    api    /swagger.yaml
    Response Should Contain    ${swagger_yaml}    openapi: 3.1.0    /v1/auctions:

Auction List Detail History And Calendar Work
    ${list_payload}=    Get JSON Response    api    /v1/auctions
    ${items}=    Get From Dictionary    ${list_payload}    data
    Should Not Be Empty    ${items}
    ${detail_payload}=    Get JSON Response    api    /v1/auctions/${LIVE_AUCTION_ID}
    ${detail}=    Get From Dictionary    ${detail_payload}    data
    Dictionary Should Contain Key    ${detail}    title
    ${history_payload}=    Get JSON Response    api    /v1/auctions/${LIVE_AUCTION_ID}/history
    ${history}=    Get From Dictionary    ${history_payload}    data
    Should Not Be Empty    ${history}
    ${calendar}=    Get Text Response    api    /v1/auctions/calendar.ics?token=${CALENDAR_TOKEN}
    Response Should Contain    ${calendar}    BEGIN:VCALENDAR    BEGIN:VEVENT

Public Content Endpoints Return Seed Data
    ${article_payload}=    Get JSON Response    api    /v1/knowledge/articles/${ARTICLE_SLUG}
    ${article}=    Get From Dictionary    ${article_payload}    data
    Dictionary Should Contain Item    ${article}    slug    ${ARTICLE_SLUG}
    ${courses_payload}=    Get JSON Response    api    /v1/courses
    ${courses}=    Get From Dictionary    ${courses_payload}    data
    Should Not Be Empty    ${courses}
    ${community_payload}=    Get JSON Response    api    /v1/community/posts
    ${community}=    Get From Dictionary    ${community_payload}    data
    Should Not Be Empty    ${community}
    ${advisors_payload}=    Get JSON Response    api    /v1/advisors
    ${advisors}=    Get From Dictionary    ${advisors_payload}    data
    Should Not Be Empty    ${advisors}
