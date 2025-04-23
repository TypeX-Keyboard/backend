# 使用官方的 Go 作为基础镜像
FROM golang:1.23 as builder


# 设置环境变量
ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.cn,direct

# 创建和设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 在容器中安装 Go 工具
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o edge-tts-go github.com/wujunwei928/edge-tts-go

# 复制项目源代码
COPY . .

# 编译 GoFrame 项目
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用更小的基础镜像
FROM debian:latest
# 安装时区设置工具
RUN apt-get update && apt-get install -y tzdata ca-certificates

# 设置时区为上海
#RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 安装 FFmpeg 和其他必要的依赖
#RUN apt-get update && apt-get install -y ffmpeg

RUN groupadd work 
RUN useradd work -m -g  work

USER work

RUN mkdir -p /home/work/app

ENV WORKDIR /home/work/app/

# 设置工作目录
WORKDIR $WORKDIR

ARG CONF=offline
# 如果需要打online环境，build时需要加 --build-arg CONF=online

# 复制配置文件到工作目录
COPY --chown=work ./manifest/ ./manifest/

# 复制静态资源文件
COPY --chown=work ./resource/ ./resource/

# 复制编译好的二进制文件到新的镜像
COPY --chown=work --from=builder /app/main ./bin/main
# COPY --chown=work --from=builder /app/edge-tts-go ./bin/edge-tts-go
# 设置环境变量
ENV PATH="${WORKDIR}/bin:${PATH}"

ENV RE=${CONF}

EXPOSE 80

# 启动应用
CMD ["./bin/main"]
