*** Settings ***
Resource         resources/common.resource
Suite Setup      Create Live Sessions

*** Test Cases ***
Home Page Loads
    ${html}=    Get Text Response    web    /
    Response Should Contain    ${html}    海關標售雷達    現狀交付    Customs Auction Platform

Auctions Page Loads
    ${html}=    Get Text Response    web    /auctions
    Response Should Contain    ${html}    全台海關標售清單    官方公告

Knowledge Community And Advisors Pages Load
    ${knowledge}=    Get Text Response    web    /knowledge/bid-form-guide
    Response Should Contain    ${knowledge}    標單
    ${community}=    Get Text Response    web    /community
    Response Should Contain    ${community}    看貨心得
    ${advisors}=    Get Text Response    web    /advisors
    Response Should Contain    ${advisors}    顧問

Manifest And Robots Are Reachable
    ${manifest}=    Get Text Response    web    /manifest.webmanifest
    Response Should Contain    ${manifest}    name    start_url
    ${robots}=    Get Text Response    web    /robots.txt
    Response Should Contain    ${robots}    User-agent: *    Sitemap:
