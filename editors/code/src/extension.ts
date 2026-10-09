import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext): void {
  const config = vscode.workspace.getConfiguration("urscript");
  const serverPath = config.get<string>("server.path", "urscript-ls");

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
