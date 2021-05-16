Package netflag implements flags for network client and server
======================================================================

サーバとクライアントは別
----------------------------------------------------------------------

    // for server
    netflag.NewServer(prefix)
    netflag.NewServerName(prefix, name)

    // for client
    netflag.NewClient(prefix)
    netflag.NewClientName(prefix, name)

プログラムの利用者の視点ではサーバとクライアントを選択する必要は
ないはずなので、フラグ名はaddress/addrで統一する。

    --address string, --addr string


networkを指定できる
----------------------------------------------------------------------

tcpとunixは入れ替えても使えそう。udpはだめそう。
オプションで指定できるのが理想的。

- デフォルトはコマンドラインでは指定できない
  - コードでネットワークを指定する
    - Flag.ListenNetwork("tcp")
    - Flag.DialNetwork("tcp")
- オプションでコマンドラインで指定可能にする
  - コマンドラインで指定する
    - Flag.Listen()
    - Flag.Dial()

    // tcpまたはunixを指定できる
    netflag.NewServer(prefix, netflag.Network("tcp", "unix"))

    // 任意のネットワークを指定できる
    netflag.NewServer(prefix, netflag.NetworkAny)

    // フラグ指定のネットワーク
    Flag.Listen()
    Flag.Dial()

    // コードで指定する
    Flag.ListenNetwork("tcp")
    Flag.DialNetwork("udp")

フラグ名はnetwork

    --network [tcp|udp|unix], --net [tcp|udp|unix]


TLSを指定できる
----------------------------------------------------------------------

tls.Configではサーバとクライアントで指定するものが異なるが、
サーバのフラグかクライアントのフラグかが決まっていれば、
同じ項目を受け取れば良さそうだ。

| S | C | 項目               |
|:--|:--|:-------------------|
| x | x | 証明書と秘密鍵     |
| x | x | MinVersion         |
| x | x | MaxVersion         |
| x |   | ClientCAs          |
|   | x | RootCAs            |

フラグは6つ。

    --tls-cert file
    --tls-cert-key file
    --tls-min-version [1.0|1.1|1.2|1.3] (default: 1.2)
    --tls-max-version [1.0|1.1|1.2|1.3] (default: 1.2)
    --tls-ca file
    --tls-skip-verify


FlagNetworkとNetworkオプション
----------------------------------------------------------------------

オプション指定の例

    Network()             // デフォルトに戻す ("tcp"固定)
    Network("*")          // 任意のネットワークの指定が必須
    Network("tcp")        // ネットワークは指定不可能
    Network("udp")        // ネットワークは指定不可能
    Network("tcp", "udp") // いずれかのネットワークを指定する
    Network("udp", "*")   // ネットワーク指定可能、デフォルトはucp
