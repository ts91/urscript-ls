import * as vscode from "vscode";
import * as path from "path";
import * as fs from "fs";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

function resolveServerPath(context: vscode.ExtensionContext): string {
  const config = vscode.workspace.getConfiguration("urscript");
  const configured = config.get<string>("server.path", "");

  // Prefer explicit user configuration.
  if (configured) {
    return configured;
  }

  // Try the bundled binary shipped with the extension.
  const binaryName =
    process.platform === "win32" ? "urscript-ls.exe" : "urscript-ls";
  const bundled = path.join(context.extensionPath, "server", binaryName);
  if (fs.existsSync(bundled)) {
    return bundled;
  }

  // Fall back to PATH.
  return "urscript-ls";
}

export function activate(context: vscode.ExtensionContext): void {
  const serverPath = resolveServerPath(context);

  const serverOptions: ServerOptions = {
    run: { command: serverPath },
    debug: { command: serverPath },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "urscript" }],
  };

  client = new LanguageClient(
    "urscript",
    "URScript Language Server",
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  return client?.stop();
}
