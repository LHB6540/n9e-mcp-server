FROM node:22-alpine

WORKDIR /app

# 允许用户在 build 时传入指定的 NPM 版本号，默认为 latest
ARG MCP_VERSION=latest

# 安装指定版本的夜莺 MCP 发布包以及网络桥接器 mcp-proxy
RUN npm config set registry https://registry.npmmirror.com && \
    npm install -g @n9e/n9e-mcp-server@${MCP_VERSION} mcp-proxy

EXPOSE 8082

# 启动代理桥接器，并透传给全局安装的夜莺 MCP 二进制
CMD ["mcp-proxy", "--port", "8082", "--", "n9e-mcp-server", "stdio"]
