*** Settings ***
Library           RequestsLibrary
Library           Collections
Library           String

*** Variables ***
${API_BASE_URL}         %{ROBOT_API_BASE_URL=https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app}
${WEB_BASE_URL}         %{ROBOT_WEB_BASE_URL=https://those-forgotten.vercel.app}
${ARTICLE_SLUG}         bid-form-guide
${CALENDAR_TOKEN}       demo

*** Keywords ***
Create Live Sessions
    Create Session    api    ${API_BASE_URL}    verify=${True}
    Create Session    web    ${WEB_BASE_URL}    verify=${True}

Get JSON Response
    [Arguments]    ${session_name}    ${path}    ${expected_status}=200
    ${response}=    GET On Session    ${session_name}    ${path}    expected_status=any
    Should Be Equal As Integers    ${response.status_code}    ${expected_status}
    ${payload}=    Evaluate    json.loads($response.text)    json
    RETURN    ${payload}

Get Text Response
    [Arguments]    ${session_name}    ${path}    ${expected_status}=200
    ${response}=    GET On Session    ${session_name}    ${path}    expected_status=any
    Should Be Equal As Integers    ${response.status_code}    ${expected_status}
    RETURN    ${response.text}

Post JSON Response
    [Arguments]    ${session_name}    ${path}    ${payload}    ${expected_status}
    ${response}=    POST On Session    ${session_name}    ${path}    json=${payload}    expected_status=any
    Should Be Equal As Integers    ${response.status_code}    ${expected_status}
    RETURN    ${response}

Prime Live Auction
    ${payload}=    Get JSON Response    api    /v1/auctions
    ${items}=    Get From Dictionary    ${payload}    data
    Should Not Be Empty    ${items}
    ${first}=    Get From List    ${items}    0
    ${auction_id}=    Get From Dictionary    ${first}    id
    Set Suite Variable    ${LIVE_AUCTION_ID}    ${auction_id}
    RETURN    ${first}

Response Should Contain
    [Arguments]    ${text}    @{expected_values}
    FOR    ${value}    IN    @{expected_values}
        Should Contain    ${text}    ${value}
    END
