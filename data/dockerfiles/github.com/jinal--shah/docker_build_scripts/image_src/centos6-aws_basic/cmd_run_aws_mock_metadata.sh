docker run -it --rm -p 80:8080 -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
        -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} jtblin/aws-mock-metadata \
        --availability-zone=a --instance-id='i-9101c68f' --hostname='antimatter-t.eurostar.com'
        --vpc-id='vpc-3cd9b759' --private-ip=10.101.34.93
