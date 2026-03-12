*** Settings ***
Resource         resources/common.resource
Suite Setup      Create Live Sessions

*** Test Cases ***
Calendar Without Token Is Rejected
    ${response}=    GET On Session    api    /v1/auctions/calendar.ics    expected_status=any
    Should Be Equal As Integers    ${response.status_code}    401

Member Endpoints Reject Anonymous Requests
    ${keyword_response}=    GET On Session    api    /v1/keyword-subscriptions    expected_status=any
    Should Be Equal As Integers    ${keyword_response.status_code}    401
    ${push_keys}=    Create Dictionary    p256dh=demo    auth=demo
    ${push_payload}=    Create Dictionary    endpoint=https://push.example.com/demo    keys=${push_keys}
    ${push_response}=    POST On Session    api    /v1/web-push-subscriptions    json=${push_payload}    expected_status=any
    Should Be Equal As Integers    ${push_response.status_code}    401
    ${community_payload}=    Create Dictionary    title=robot post    body=robot test body    office=臺北關
    ${community_response}=    POST On Session    api    /v1/community/posts    json=${community_payload}    expected_status=any
    Should Be Equal As Integers    ${community_response.status_code}    401
    ${checkout_payload}=    Create Dictionary    kind=membership    plan_code=pro-monthly
    ${checkout_response}=    POST On Session    api    /v1/stripe/checkout    json=${checkout_payload}    expected_status=any
    Should Be Equal As Integers    ${checkout_response.status_code}    401

Admin Endpoints Reject Anonymous Requests
    ${reports_response}=    GET On Session    api    /v1/admin/community-reports    expected_status=any
    Should Be Equal As Integers    ${reports_response.status_code}    403
    ${crawler_response}=    GET On Session    api    /v1/admin/crawler-status    expected_status=any
    Should Be Equal As Integers    ${crawler_response.status_code}    403
    ${advisor_response}=    GET On Session    api    /v1/admin/advisor-leads    expected_status=any
    Should Be Equal As Integers    ${advisor_response.status_code}    403
