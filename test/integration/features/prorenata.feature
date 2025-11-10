@prorenata
Feature: Features added on Prorenata fork

  Scenario: POST /forms/gspreview (pdf -> png)
    Given I have a default Gotenberg container
    When I make a "POST" request to Gotenberg at the "/forms/gspreview" endpoint with the following form data and header(s):
      | files                     | testdata/page_1.pdf  | file   |
      | Gotenberg-Output-Filename | foo                  | header |
    Then the response status code should be 200
    Then the response header "Content-Type" should be "image/png"
    Then there should be the following file(s) in the response:
      | foo.png |

  Scenario: POST /forms/gspreview (pdf -> png with with arguments)
    Given I have a default Gotenberg container
    When I make a "POST" request to Gotenberg at the "/forms/gspreview" endpoint with the following form data and header(s):
      | files                     | testdata/page_1.pdf  | file   |
      | Gotenberg-Output-Filename | foo                  | header |
      | xsize                     | 600                  | field  |
      | outputFormat              | png                  | field  |
    Then the response status code should be 200
    Then the response header "Content-Type" should be "image/png"
    Then there should be the following file(s) in the response:
      | foo.png |
    
  Scenario: POST /forms/gspreview (image -> pdf)
    Given I have a default Gotenberg container
    When I make a "POST" request to Gotenberg at the "/forms/gspreview" endpoint with the following form data and header(s):
      | files                     | testdata/test.png  | file   |
      | Gotenberg-Output-Filename | foo                | header |
      | outputFormat              | pdf                | field  |
    Then the response status code should be 200
    Then the response header "Content-Type" should be "application/pdf"
    Then there should be the following file(s) in the response:
      | foo.pdf |

  Scenario: POST /forms/libreoffice/convert (rtf -> plain text)
    Given I have a default Gotenberg container
    When I make a "POST" request to Gotenberg at the "/forms/libreoffice/convert" endpoint with the following form data and header(s):
      | files                     | testdata/test.rtf  | file   |
      | Gotenberg-Output-Filename | foo                | header |
      | outputFormat              | text               | field  |
    Then the response status code should be 200
    Then the response header "Content-Type" should be "text/plain; charset=utf-8"
    Then there should be the following file(s) in the response:
      | foo.txt |
