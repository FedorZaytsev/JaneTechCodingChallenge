# Solution for Jane Technologies coding challenge

## Techincal details
This is an interesting problem because of the fact that the data should be accessable by 2 different keys - campaign id and keywords. In this solution I used 2 simple hash maps to store the data, one is keyed by campaign id and used to increment add impressions quick and second one is more like an index and keyed by target keywords to do fast access to the matching campaigns. In memory DB is guarded by database connector class and provide a simplish interface with the idea that in a real world scenario this would be replaced with a standalone DB. The choice of the standalone DB requires knowledge of many parts missing in the coding challenge, both SQL and NoSQL are viable in my opinion, however I am inclining towards NoSQL for better performance.

The most interesting handler here is /addecission and I will talk more about it. Because we don't know exact underlying DB I decided to go less optimal but safer way and for each keyword I fetch all matching campaigns regardless of their CPM. This is suboptimal and if we assume that we store everything in memory could be greatly optimized (see connector.go FetchMatchingCampaign) by assuming that all the data could be sorted.

## How to run
Just run `./run.sh` it would build and run the server. Config is located in the `./configs` folder

## Failing tests
If I understand the task correctly, provided test cases are not fully correct. I have problems with #6 and #10

### Test #6
In first two test requests you create two exact same campaigns with different CPM. The important detail here that the target keywords are the same too. Then in request #5 you record an impression to one of the two campaigns, however the second campaign is still active and has the same target keywords. So request #6 should return second campaign instead of null:
6. `curl -i --header 'Content-Type: application/json' http://localhost:8000/addecision --data-raw '{"keywords":["iphone","5G"]}' `
200 OK. 
response body is empty. 
There is no matching campaign because the campaign 1001 has reached its max impressions.

### Test #10
Test #10 is identical to the test #9, but have different results. I assume you need to change campaign id