FROM mcr.microsoft.com/dotnet/sdk:8.0-preview-cbl-mariner2.0

RUN cd /tmp \
    && curl -LO https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 \
    && mv ./jq-linux-amd64 /usr/bin/jq \
    && chmod +x /usr/bin/jq

WORKDIR /build/appinspector

COPY ./ApplicationInspector/ ./src

WORKDIR /build/appinspector/src

RUN set -o errexit -o nounset && \
    dotnet build -c Release && \
    ln -s /build/appinspector/src/AppInspector.CLI/bin/Release/net8.0/ApplicationInspector.CLI /usr/bin/appinspector

WORKDIR /build/OSSGadget

COPY ./OSSGadget/src ./src

WORKDIR /build/OSSGadget/src

RUN set -o errexit -o nounset \
    && dotnet build -c Release \
    && ln -s /build/OSSGadget/src/oss-download/bin/Release/net6.0/oss-download /usr/bin/oss-download \
    && ln -s /build/OSSGadget/src/oss-find-source/bin/Release/net6.0/oss-find-source /usr/bin/oss-find-source \
    && ln -s /build/OSSGadget/src/oss-metadata/bin/Release/net6.0/oss-metadata /usr/bin/oss-metadata

WORKDIR /app

COPY . .

ENTRYPOINT ["sh", "/app/docker-entrypoint.sh"]
