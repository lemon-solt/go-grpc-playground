FROM golang:1.20
COPY . .
RUN apt-get update && apt-get install -y bash
CMD ["bash"]