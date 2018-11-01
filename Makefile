b:
	rm sendemail && rm sendemail.zip && \
	env GOOS=linux go build -v -ldflags '-d -s -w' \
	-a -tags netgo -installsuffix netgo -o sendemail main.go &&\
	zip sendemail.zip ./sendemail 