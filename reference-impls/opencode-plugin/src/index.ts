import type { Plugin } from '@opencode-ai/plugin';
import { GrpcTransport } from "@protobuf-ts/grpc-transport";
import { ChannelCredentials } from "@grpc/grpc-js";
import { Application } from './modelhawk/application';
import { NotifyServiceClient, } from "./modelhawk/notify_service.client";
import { InfoServiceClient, } from "./modelhawk/info_service.client";
import { Timestamp } from "./modelhawk/google/protobuf/timestamp";

export const ModelHawkClient: Plugin = async ({ client, directory }) => {
    await client.app.log({
        body: {
            service: "model-hawk-client",
            level: "info",
            message: "started",
        },
    });

    const app: Application = { value: "my-opencode-session" };
    const reportedTools = new Set<string>();
    const port = 50051;

    // Create a gRPC transport pointing at your server
    const transport = new GrpcTransport({
        host: `localhost:${port}`,
        channelCredentials: ChannelCredentials.createInsecure(),
    });

    // make clients
    const notifyClient = new NotifyServiceClient(transport);
    const infoClient = new InfoServiceClient(transport);

    return {
        "tool.execute.before": async (input, output) => {
            if (!reportedTools.has(input.tool)) {
                // tell the ModelHawk server about the tool
                const resp = await infoClient.giveToolInfo({
                    app: app,
                    name: input.tool,
                    args: [],
                });
                reportedTools.add(input.tool);
            }

            switch (input.tool) {
                case "read":
                    await notifyClient.calledTool({
                        app: app,
                        name: "read",
                        args: [],
                    });
            }
        },
    };
};