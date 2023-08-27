# README
## AppsyncからGoAPIへの通信
AWS AppSyncは、GraphQLサービスであり、リゾルバを通じてデータソースとの間でデータの取得や変更を行います。AppSyncを設定する際に、リゾルバを使って特定のGraphQLクエリやミューテーションがデータソース（例：RDS、DynamoDB、Lambda、HTTPデータソースなど）にどのように対応するかを定義します。

AppSyncがHTTPやRESTful API（この場合、GoAPI）と通信する場合、以下のようなリクエストが考えられます：

1. **POSTリクエスト**: GraphQLの特性上、AppSyncからのリクエストは通常、HTTP POSTリクエストとして送信されます。

2. **エンドポイント**: あなたがAppSyncのリゾルバで指定したエンドポイントにリクエストが送信されます。例えば、GoAPIの`/messages`エンドポイントにデータを取得するクエリをマッピングする場合、そのエンドポイントにリクエストが送信されます。

3. **ヘッダー**: AppSyncからのリクエストには、必要に応じてヘッダーが含まれることがあります。これには認証情報やカスタムヘッダーなどが含まれることがあります。

4. **リクエストボディ**: GraphQLクエリやミューテーションの情報がリクエストボディにJSON形式でエンコードされています。リゾルバで指定したマッピングテンプレートに基づいて、このボディは変換され、APIが期待する形式になります。

具体的な例：

AppSyncのGraphQLミューテーションが以下のような場合：

```graphql
mutation sendMessage($content: String!) {
    sendMessage(content: $content) {
        id
        content
    }
}
```

そして、このミューテーションがGoAPIの`/send`エンドポイントにマッピングされていると仮定すると、AppSyncからGoAPIには以下のようなHTTP POSTリクエストが送信されます：

```json
{
  "content": "Hello from AppSync!"
}
```

しかし、実際のリクエストの形式は、AppSyncで設定したリゾルバのマッピングテンプレートによって異なります。リゾルバのマッピングテンプレートを使用して、AppSyncからの入力をAPIが期待する形式に変換することができます。