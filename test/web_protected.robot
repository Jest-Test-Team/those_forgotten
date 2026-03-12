*** Settings ***
Resource         resources/common.resource
Suite Setup      Create Live Sessions

*** Test Cases ***
Member Page Must Not Crash
    ${response}=    GET On Session    web    /member    expected_status=any
    Should Be True    ${response.status_code} < 500

Admin Page Must Not Crash
    ${response}=    GET On Session    web    /admin    expected_status=any
    Should Be True    ${response.status_code} < 500
