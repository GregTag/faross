FROM mcr.microsoft.com/dotnet/sdk:8.0-preview-cbl-mariner2.0

RUN cd /tmp \
    && curl -LO https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 \
    && mv ./jq-linux-amd64 /usr/bin/jq \
    && chmod +x /usr/bin/jq

COPY ./OSSGadget/ /app/

WORKDIR /app/src

RUN set -o errexit -o nounset \
    && dotnet build -c Release \
    && ln -s /app/src/oss-characteristics/bin/Release/net6.0/oss-characteristic /usr/bin/oss-characteristic \
    && ln -s /app/src/oss-defog/bin/Release/net6.0/oss-defog /usr/bin/oss-defog \
    && ln -s /app/src/oss-detect-backdoor/bin/Release/net6.0/oss-detect-backdoor /usr/bin/oss-detect-backdoor \
    && ln -s /app/src/oss-detect-cryptography/bin/Release/net6.0/oss-detect-cryptography /usr/bin/oss-detect-cryptography \
    && ln -s /app/src/oss-diff/bin/Release/net6.0/oss-diff /usr/bin/oss-diff \
    && ln -s /app/src/oss-download/bin/Release/net6.0/oss-download /usr/bin/oss-download \
    && ln -s /app/src/oss-find-domain-squats/bin/Release/net6.0/oss-find-domain-squats /usr/bin/oss-find-domain-squats \
    && ln -s /app/src/oss-find-source/bin/Release/net6.0/oss-find-source /usr/bin/oss-find-source \
    && ln -s /app/src/oss-find-squats/bin/Release/net6.0/oss-find-squats /usr/bin/oss-find-squats \
    && ln -s /app/src/oss-health/bin/Release/net6.0/oss-health /usr/bin/oss-health \
    && ln -s /app/src/oss-metadata/bin/Release/net6.0/oss-metadata /usr/bin/oss-metadata \
    && ln -s /app/src/oss-risk-calculator/bin/Release/net6.0/oss-risk-calculator /usr/bin/oss-risk-calculator

COPY ./docker-entrypoint.sh .

ENTRYPOINT ["sh", "/app/src/docker-entrypoint.sh"]
