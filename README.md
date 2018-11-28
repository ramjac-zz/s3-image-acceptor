# S3 Image Acceptor

This is a lambda function that accepts a posted image and saves it to AWS S3 in a bucket determined by the function's configuration. Returns the full, schemaless URL of the image in S3 on success.

In addition to the image, these optional inputs will also be read in the post:

* expiration - An integer for the number of hours until the image should expire. Default is to never expire.
* name - The name the image file will have in S3.

Example HTTP POST

