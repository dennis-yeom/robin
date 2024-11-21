# robin
Robin performs ETL on data.

### Function
Ticker will trigger us to check for a message.
If message is found, this indicates that batman has detected a change in the files.
We go to the bucket to access the files.
For each file, we check the file and see if it is in our Redis DB.
If it is not in our Redis DB, that means we have not downloaded that file yet.
Store the object name and object version into Redis.
Download that file into Mongo DB.
After doing that for all the files, wait for next ticker to trigger.

# SQS
We will listen for a message on an AWS message queue. We will be checking the queue periodically using a ticker.

# Linode 
Linode access for access to buckets. Our data will be stored in json files.

# MongoDB
Possble database to use?

# Laravel
PHP-based framework. Serve a web interface for managing and viewing the processed data. Provide API endpoints for external applications to interact with the ETL pipeline.

# Cobra
For CLI

# Viper
For environment variables
