const path = require("path");

module.exports = {
    entry: "./src/main.ts", // Update the path to point to the src folder
    output: {
        filename: "main.js",
        path: path.resolve(__dirname)
    },
    resolve: {
        extensions: [".ts", ".js"],
        alias: {
            "@src": path.resolve(__dirname, "src") // Create an alias for the src folder
        }
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: "ts-loader",
                exclude: /node_modules/
            }
        ]
    },
    watchOptions: {
        ignored: /node_modules/
    }
};
