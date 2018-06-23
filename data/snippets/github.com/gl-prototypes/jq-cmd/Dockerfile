FROM alpine:3.6
RUN apk add -U jq
COPY ./jqserve /jqserve
RUN chmod +x /jqserve
EXPOSE 2222
ENTRYPOINT ["/jqserve"]
