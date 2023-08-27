# ページング
## カーソルベースのページング
GraphQLでよく採用されるページングの方法は、「カーソルベースのページング」として知られています。この方法では、カーソル（一意のポインタ）を使用してデータの特定の位置を参照します。

1. **カーソル**: 通常、エンコードされた文字列として提供されるカーソルは、データセットの特定の位置を指します。カーソルは一般に、データの特定のアイテム（たとえばIDやタイムスタンプ）に基づいて生成されます。

2. **forward pagination**と**backward pagination**: クライアントは`after`パラメータを使用して特定のカーソルの後のアイテムを取得するか、`before`パラメータを使用して特定のカーソルの前のアイテムを取得します。

3. **エンドポイントの一般化**: 同じエンドポイントを使用して、前方または後方のページングを実行できます。

4. **セキュリティ**: カーソルはエンコードされた文字列として提供されるため、内部的なデータの構造やIDが外部に露出するリスクが低くなります。

GraphQLの一般的なページングは、以下のようになります：

```graphql
{
  messages(first: 10, after: "YXJyYXljb25uZWN0aW9uOjE=") {
    edges {
      node {
        id
        content
      }
      cursor
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

ここでの応答は、`edges`と`pageInfo`の2つの主要な部分で構成されます。`edges`は実際のデータのリストであり、各`edge`は`node`（実際のデータ）と`cursor`（そのデータのカーソル）を持っています。`pageInfo`はページングに関する情報（次のページがあるかどうか、次のページのカーソルなど）を提供します。

このカーソルベースのページング方法は、特にデータセットが大きい場合やデータが頻繁に変更される場合に、オフセットベースのページングよりもパフォーマンスが優れていると言われています。


## IDをoffsetとしてページングする方法
IDをオフセットとして使用するページングは、データベースのクエリで `WHERE` 句と `LIMIT` 句を使用して実現します。以下はその方法と、GoとPostgreSQLを使用した実例を示しています。

### 基本的なアイディア:

1. クライアントが最初のページをリクエストするときには、オフセットのIDを指定しない。
2. クライアントが次のページをリクエストするときには、前のページの最後のアイテムのIDをオフセットとして指定する。
3. サーバーは、そのオフセットIDを基にして次のページのアイテムを取得する。

### GoとPostgreSQLを使用した実例:

1. **データベースのクエリ**:

```sql
-- 最初のページを取得する場合
SELECT * FROM messages ORDER BY id ASC LIMIT 10;

-- 2ページ目以降を取得する場合 (例: オフセットIDが20の場合)
SELECT * FROM messages WHERE id > 20 ORDER BY id ASC LIMIT 10;
```

2. **Goのコード**:

まず、必要なパッケージをインポートします:

```go
import (
    "database/sql"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)
```

次に、エンドポイントを実装します:

```go
func getMessages(c *gin.Context) {
    var messages []Message
    var err error

    limit := 10  // 1ページあたりのアイテム数
    offsetID := c.DefaultQuery("offsetID", "0")

    lastID, _ := strconv.Atoi(offsetID)

    // クエリの実行
    rows, err := db.Query("SELECT * FROM messages WHERE id > $1 ORDER BY id ASC LIMIT $2", lastID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.Content); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        messages = append(messages, msg)
    }

    c.JSON(http.StatusOK, messages)
}
```

このコードでは、`offsetID` クエリパラメータを使用して、クライアントが指定したIDを基にページングを実行します。最初のページをリクエストする場合、クライアントはこのパラメータを指定しないため、デフォルトの `0` が使用されます。

この方法を使用すると、オフセットをIDとして簡単に使用できますが、データベースのレコードが削除された場合など、ページングの間にデータが変更されると、期待した結果が得られない可能性があります。


### idをoffsetとしてページングする方法に加えて日付を利用してソート
IDをオフセットとして使用しつつ、日付でソートするページングは、`WHERE`句、`ORDER BY`句、および`LIMIT`句を適切に組み合わせることで実現します。

以下に、GoとPostgreSQLを使用した例を示します。ここでは、`messages`テーブルに`created_at`という日付のカラムがあると仮定します。

### データベースのクエリ:

1. **最初のページを取得**:
```sql
SELECT * FROM messages ORDER BY created_at DESC, id DESC LIMIT 10;
```

2. **2ページ目以降を取得** (例: `created_at`が特定の日付で、`id`が20の場合):
```sql
SELECT * FROM messages 
WHERE (created_at, id) < ('2023-05-01', 20) 
ORDER BY created_at DESC, id DESC 
LIMIT 10;
```

### Goのコード:

```go
import (
    "database/sql"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func getMessages(c *gin.Context) {
    var messages []Message
    var err error

    limit := 10  // 1ページあたりのアイテム数
    offsetDate := c.DefaultQuery("offsetDate", "")
    offsetID := c.DefaultQuery("offsetID", "0")

    if offsetDate == "" {
        // 最初のページの取得
        rows, err := db.Query("SELECT * FROM messages ORDER BY created_at DESC, id DESC LIMIT $1", limit)
        // 以下の処理は変わりません...
    } else {
        parsedDate, err := time.Parse("2006-01-02", offsetDate)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use 'YYYY-MM-DD'."})
            return
        }

        lastID, _ := strconv.Atoi(offsetID)

        // クエリの実行
        rows, err := db.Query(
            "SELECT * FROM messages WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC, id DESC LIMIT $3",
            parsedDate, lastID, limit,
        )
        // 以下の処理は変わりません...
    }

    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        messages = append(messages, msg)
    }

    c.JSON(http.StatusOK, messages)
}
```

この方法では、日付（`created_at`）で降順ソートし、同じ日付の場合はID（`id`）でさらに降順ソートします。クライアントは前のページの最後のメッセージの日付とIDをオフセットとして次のページをリクエストします。