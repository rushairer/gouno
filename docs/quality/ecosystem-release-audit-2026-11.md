# Gouno Ecosystem Release Audit — 2026-11

审计日期：2026-11-05（Asia/Shanghai）。本报告区分已发布、仅准备完成与未验收，不以编译、签名步骤或绿色自动化替代完整发布验收。

## 结果与范围

基础创建链 Core v1.4.0 → Template v1.4.0 → CLI v1.3.0 已公开；Agent Demo v0.3.0 独立发布。Gosso v1.7.0 与文档 v1.3.1 也已正式发布并完成各自产物验收。GOSSO Admin 与 Blog 仍受认证、部署验收及上游门禁约束，不得创建正式候选 Tag。

已读取七个现有仓库及独立克隆的 Gouno UI 的 AGENTS.md；Agent Demo 仓库没有 AGENTS.md。所有原工作区最初干净，原有未提交工作没有被覆盖。未 force push、删除或覆盖任何 Tag。Blog 原克隆普通 fetch 拒绝多个历史 Tag 冲突（v1.0.0、v1.2.x、部分 v1.3.x/v1.4.x/v1.5.x）；保留原对象，通过独立工作副本和公开 peeled ref 核对 v1.5.26。未据此宣称远端曾改写 Tag；冲突原因仍待历史核查。

## 版本与不可变提交

| 仓库 | 开始时正式基线 | 本次目标 | 实际正式 Tag peeled SHA | 状态 |
| --- | --- | --- | --- | --- |
| gouno | v1.4.0 已有 Tag/Go Module，缺 Release | v1.4.0 | `20efba0af181cb3e03c5bcb199576bfcdec1bd7f` | 原 Tag 未改，补建 GitHub Release |
| gouno-template | v1.3.0 | v1.4.0 | `9f0cba9c8a6be6331f48596c83663efbec5490a2` | annotated Tag 与 Release 已公开 |
| gouno-cli | v1.2.1 | v1.3.0 | `087c00a523a6ec8c512d5eb9c65f657f6be8d5fb` | annotated Tag、Release、Proxy 已验证 |
| gosso | v1.6.1 | v1.7.0 | `5bd45c8f9946662274117c4644b31ce90a003464` | annotated Tag、Release、镜像、签名、SBOM 与 Helm OCI 验收通过 |
| gosso-admin | v0.7.24 | v0.7.25 | 未创建 | 阻塞，未发布 |
| gouno-blog | v1.5.26 | 候选 v1.6.0 | 未创建 | 阻塞，未发布 |
| gouno-ui | npm/Tag v0.4.9 | 保持 v0.4.9 | `4495147ae0b9807d2e5c9137187cfdc3e906cceb` | 无新包发布 |
| gouno-doc | v1.3.0 | v1.3.1 | `fad499f7f28ff2e86d03d5c26ae6de127222d915` | 中英文矩阵、CI、annotated Tag 与 Release 完成 |
| gouno-agent-demo | CHANGELOG 声明 v0.2.0，远端无 Tag/Release | v0.3.0 | `ca372fc72379dbf63e6794f5d4e65af7c88e8ba7` | 双版本最终 CI 后独立发布 |

原正式基线的远端 peeled SHA（便于重现差异审核）：

| 基线 | SHA |
| --- | --- |
| gouno v1.4.0 | `20efba0af181cb3e03c5bcb199576bfcdec1bd7f` |
| gouno-template v1.3.0 | `c7d09d3393d460f8ebe8823f6c6175f42e83c0d8` |
| gouno-cli v1.2.1 | `9892c689d2e383593ee92d65f94303a870d376f0` |
| gosso v1.6.1 | `0ec136f67b2fa165abbf77fedb5fdd50b5ebe69f` |
| gosso-admin v0.7.24 | `974ddfce7c4efe5abe7edfa7b4e597e9258573d3` |
| gouno-blog v1.5.26 | `107ee2e28eb96bda1409926fdb324550bb1bae98` |
| gouno-ui v0.4.9 | `4495147ae0b9807d2e5c9137187cfdc3e906cceb` |
| gouno-doc v1.3.0 | `55d5ee2758457ffa3d7a550ed1176824e8e0f255` |

公开 Release：[Gosso v1.7.0](https://github.com/rushairer/gosso/releases/tag/v1.7.0)、[Doc v1.3.1](https://github.com/rushairer/gouno-doc/releases/tag/v1.3.1)、[Core v1.4.0](https://github.com/rushairer/gouno/releases/tag/v1.4.0)、[Template v1.4.0](https://github.com/rushairer/gouno-template/releases/tag/v1.4.0)、[CLI v1.3.0](https://github.com/rushairer/gouno-cli/releases/tag/v1.3.0)、[Agent Demo v0.3.0](https://github.com/rushairer/gouno-agent-demo/releases/tag/v0.3.0)。

原正式下游：[Gosso v1.6.1](https://github.com/rushairer/gosso/releases/tag/v1.6.1)、[Admin v0.7.24](https://github.com/rushairer/gosso-admin/releases/tag/v0.7.24)、[Blog v1.5.26](https://github.com/rushairer/gouno-blog/releases/tag/v1.5.26)。这些链接不是本次候选发布证明。

新建 annotated Tag 对象：Template `75a4d30b3a3eca0dcda339daf1173eaac4c6ef31`；CLI `b898d5e921fb9a001e03a65b1b880293794fc606`；Gosso `15d529590b119621a51d624d192f489591c12bf9`；Doc `b3797d0260806d8b8ea84a72d8b223070116c54e`；Demo `18bf25e95bccc3dda0e6bae0cb680abff7bfc748`。Core 原 annotated Tag 对象 `1ccd3bed5de725d90cd318a7f4799da2b0eec079` 保持不变。Tag 对象 SHA 与 peeled commit SHA 分别核对。

显式设置官方 `GOPROXY=https://proxy.golang.org`、`GOSUMDB=sum.golang.org` 并使用全新 module cache，重新下载 Core、CLI、Gosso 三个版本，均退出 0；并再次核对 npm tarball 与两个消费者 lockfile 的 version/resolved/integrity。

本地验收工具：Go 1.26.10 / 1.27.3、Node 24.18.0、Docker Engine 29.9.2、Cosign 3.1.4、Helm 4.2.3、Trivy 0.76.1。Cosign 二进制 SHA256 与官方 GitHub asset digest 一致，Helm archive 与官方 SHA256 文件一致。

## Go Module、npm 与兼容性

| 产物 | 公开解析与校验 |
| --- | --- |
| github.com/rushairer/gouno v1.4.0 | Proxy Origin Hash=`20efba0…`；module `h1:cCs2eGv6XCJMN4q2RFaFPCL1Ip0PY21RSHp3mT4tAYY=`；go.mod `h1:S2DmlM4BlX3BLoeRw7s66aABDJDexP4Gj/xov28MmvA=` |
| github.com/rushairer/gouno-cli v1.3.0 | Proxy Origin Hash=`087c00a…`；module `h1:dL7sChra+i9zgpNxmdHFc/tYlBGhDHqfhHC0YWsyLD0=`；go.mod `h1:UVjoYQk7M5sC7ziidDTT3zrG/KqzYBL8pSdM1kH1Pls=` |
| github.com/rushairer/gosso v1.7.0 | Proxy Origin Hash=`5bd45c8…`；module `h1:ZdnBcJycBnpKmpQ1hFLnOhzCj4DB2f6JYmhbGQY10s0=`；go.mod `h1:5XTaKie4g0Vwjo0r95irL9Ph2MnXHA+t1EbriqRM2+w=` |
| gouno-template v1.4.0 | Git Tag 验收；根 go.mod 是 `module {{.ModulePath}}` 模板，不向 Go Proxy 发布 |
| @gouno/ui 0.4.9 | npm registry 元数据及实际 tarball SHA512 一致：`sha512-+mOkw+VL/whS2bBK8mQlU6kMHjL5AxleOE1qO2rK3xD58KZsQ6IyO+0S/UciKklvrJIH+qz/LhMUrJ04XxWjlQ==` |
| @gosso/client 0.9.3 | npm registry 元数据及实际 tarball SHA512 一致：`sha512-weH+uD44qTxV6rGnvlaXXj9VHkulCiq6vk6NoHgd9P19eH8VL9e9g0G7aJOHcR+0nrJaD2yuYUXgAEUpI/F/Pg==` |
| Agent Demo v0.3.0 | Git Tag/Release；Proxy `.info` 可解析该 SHA。其 go.mod 的应用路径是 `gouno-agent-demo`，不将 `.info` 当作可按 GitHub 路径导入的库兼容承诺 |

依赖矩阵：

| 消费者 | Core | UI / Client | 工具链 | 发布边界 |
| --- | --- | --- | --- | --- |
| Template v1.4.0 | v1.4.0，go.sum 与公开校验和一致 | 不适用 | Go >=1.26；CI 1.26.10、1.27.3 | 已发布 |
| CLI v1.3.0 | Bootstrap 契约；不强加项目架构 | 固定 Template v1.4.0 与 live main 都验证 | Go >=1.26；CI 双系列 | 已发布 |
| Gosso v1.7.0 | v1.4.0 | 不适用 | Go 1.26.x；lint >=2.15.0 | 已发布，正式产物验收通过 |
| Admin v0.7.25 候选 | Seed 独立依赖 | UI 0.4.9 / Client 0.9.3 | Node >=24.15 <25 | 全局 Client 治理需要 Gosso v1.7.0；v1.6.1 的旧 smoke 不能认证新 API |
| Blog v1.6.0 候选 | v1.4.0 | UI 0.4.9 / Client 0.9.3 | Go 1.27.3；Node 24 LTS | 已验证正式 Gosso v1.7.0 精确 digest 的 BFF 兼容性；该证据属于候选 main，不是已发布 Blog Tag |
| Demo v0.3.0 | v1.4.0 | 不适用 | Go >=1.26；双系列 CI | 独立发布，不阻塞主链 |

## CI 与安全验收

| 仓库 / 精确提交 | 证据与实际结论 |
| --- | --- |
| Core Tag `20efba0…` | [CI 41813309775](https://github.com/rushairer/gouno/actions/runs/41813309775) 双 Go 系列、lint 完成成功；main 后续提交 `3171d57…` 的 latest-patch CI 也成功，但 Release 仍使用原 Tag 的 CHANGELOG 与 SHA |
| Template Tag `9f0cba9…` | [CI 41817614153](https://github.com/rushairer/gouno-template/actions/runs/41817614153)：Go 1.26.10 / 1.27.3，tidy 无差异、单项 service 和组合 suite、Race、健康路由、Vet、Lint、govulncheck 全过。扫描为 0 reachable，另报告 3 个未调用模块 advisory，并未隐去 |
| CLI Tag `087c00a…` | [CI 42010113946](https://github.com/rushairer/gouno-cli/actions/runs/42010113946) 六个 job 全成功；live default 和 pinned v1.4.0 均自动 tidy、保留 Codegen、单项/组合生成、Race/Vet。去掉 pinned 测试原先的 skip-tidy，扩大为真正生成项目验收；本地 public-tag Bootstrap 另做 `.gouno` 逐字节比较、Race、Vet |
| Demo Tag `ca372fc…` | [CI 42011180831](https://github.com/rushairer/gouno-agent-demo/actions/runs/42011180831) Go 1.26.x/1.27.x quality、lint、Conventional PR 全成功。本地 1.26.10/1.27.3 Race/Vet、module verify、govulncheck 成功；Agent、HTTP JSON/SSE 和 OpenAI/Anthropic mock HTTP Provider 测试通过，无生产 Provider 凭据调用 |
| Gosso Tag `5bd45c8…` | [CI 42011835603](https://github.com/rushairer/gosso/actions/runs/42011835603) 完成成功：Race、Vet、Lint、gosec、govulncheck、PostgreSQL/Redis HTTP integration、migration、Helm、Docker 与镜像扫描；coverage=70.0%，critical package ratchet 保持。PR-only Conventional Commits 在 main push 合理跳过；[Release Workflow 42013419089](https://github.com/rushairer/gosso/actions/runs/42013419089) 亦完成成功 |
| Admin main `51fe744…` | [CI](https://github.com/rushairer/gosso-admin/actions/runs/42016212154)、[Images](https://github.com/rushairer/gosso-admin/actions/runs/42016212304)、[Security](https://github.com/rushairer/gosso-admin/actions/runs/42016212099) 成功。实际代码提交 `64f4c43…` 的 [compatibility](https://github.com/rushairer/gosso-admin/actions/runs/42015649627) 成功，覆盖正式 Gosso v1.7.0 的全局 Client CRUD/rotate、Seed、Cookie 与 PKCE。后续仅更新发布清单；未改 Product UI，之前 `0feeab7…` 的 [browser](https://github.com/rushairer/gosso-admin/actions/runs/41811416141) 和 [parity](https://github.com/rushairer/gosso-admin/actions/runs/41811416037) 仍只作支持证据，不续期人工认证 |
| Blog main `53d4a7c…` | [CI](https://github.com/rushairer/gouno-blog/actions/runs/42016283335)、[Images](https://github.com/rushairer/gouno-blog/actions/runs/42016283318)、[BFF](https://github.com/rushairer/gouno-blog/actions/runs/42016283306) 全部成功。DB artifact 13409786334 下载 zip 的实际 SHA256=`156bafca883a619fbfb67f0b7632587d87bd8a087fba51d23f396320097ec056`，与 GitHub digest 一致；独立检查 summary 为 21 packages / 115 required tests，failed/skipped/missing=0，cleanup=0。三项失败探针按预期退出 1；PR #414 的 dependency review 亦通过 |
| UI main `b706b6c…` | [CI](https://github.com/rushairer/gouno-ui/actions/runs/41805317263)、[Pages](https://github.com/rushairer/gouno-ui/actions/runs/41805317403)、[Golden](https://github.com/rushairer/gouno-ui/actions/runs/41805317297)、[Admin parity](https://github.com/rushairer/gouno-ui/actions/runs/41805317291)、[Blog parity](https://github.com/rushairer/gouno-ui/actions/runs/41805317271) 均成功；gh-pages `883fa206…` 部署提交明确引用 main `b706b6c…` |
| Doc Tag `fad499f…` | 本地 `scripts/check-docs.sh` 与 [CI 42015368871](https://github.com/rushairer/gouno-doc/actions/runs/42015368871) 完成成功；中英文矩阵只记录真实发布产物，不把候选提交的兼容证据附会到旧 Admin/Blog Tag |

Template 和 CLI 的发布后证据分别提交到 main `01f227f…`、`a3930d1…`，各自新 CI 成功；没有移动已发布 Tag。UI 从 v0.4.9 到 main 的 `src/`、package.json、package-lock.json diff 为空；变化归属 Showcase、文档、测试与 CI，没有发布 v0.4.10 的包运行时理由。

## 下游审核结论与限制

Gosso：提交 `4e964ede…` 添加全局管理员 API，满足 backward-compatible minor v1.7.0 条件。实际路由在 `/api/v1/admin/oauth2/clients`，统一经过限流、JWT、AdminRequired；中间件要求 user_session principal、admin scope 和 admin role，机器及委托 principal 不因此获得管理权限。普通路由继续校验 owner；更新保留同事务读写和 optimistic locking；ClientSecretHash 为 json:"-"。Admin clientService 的五类请求与这些路由匹配。Chart version/appVersion 与 dated CHANGELOG 已对齐 1.7.0。

Gosso 原 standalone gosec `-no-fail` 不能证明安全通过。首次独立扫描 21 项：19 项逐条复核为已有范围校验、公开常量、动态 production Cookie 或部署控制路径，使用有理由的行级标注；两处部署工具日志移除不可信变量。另复核 lint 内嵌 scanner 暴露的 20 个测试 fixture 发现；移除六类全局 lint 排除。最终本地 standalone gosec 0 unsuppressed、完整 lint 0 issues、关键包 Race 成功。处置依据保存在 Gosso `docs/quality/release-security-disposition-2026-11.md`，并没有禁用 detector 类别或隐藏 reachable govulncheck findings。

Admin：v0.7.24 到 main 的差异不只是旧 CHANGELOG 所述 CSS。已重写尚未发布的 v0.7.25 条目，记录全局 Client 治理、Canonical UI 0.4.9、Client 0.9.3、Node 24 和安全/交互变化。[PR #88](https://github.com/rushairer/gosso-admin/pull/88) 在全套 CI/安全门禁及正式 v1.7.0 镜像的全局治理 smoke 通过后合并为 `64f4c437898eec66962cb966c840588eb331c9e1`，后续清单证据提交为 `51fe7448e3ff526e961630133443221507e9633f`。认证唯一权威 JSON 当前仍是 `needs-manual-recertification`，即使 freshness 脚本输出 valid/current，也不能解释为人工认证通过。发布清单要求真实 staging 的生产配置、可变更 smoke、备份恢复、回滚等，当前没有已识别授权 staging 位置/操作资料，已请求这些信息；未发送或索取秘密。只读 Docker inventory 没有 Gouno staging 栈，现有无关容器未被触碰；读取 docs/OPERATIONS.md 确认备份必须有真实恢复演练。

Blog：初始 `v1.5.26..18e301ef` 为 780 files、90387 insertions / 37712 deletions。新外部 capability gateway、API Access 管理、Provider 配置和守卫式恢复能力支持 v1.6.0 候选，不能按“仅补丁”机械选 v1.5.27。已整理所有 Unreleased 分类，保留历史发布条目；[PR #412](https://github.com/rushairer/gouno-blog/pull/412) 通过完整 CI 后合并为 `35138fed43ffc8ae2d0b887912a8ff867b854aa9`。

抽查生产实现与已验收回归：Agent `commitApprovalEffect` 在同一 dbtx 中持久化业务效果和 approved→executed；Page/Post 保留最终 SQL CAS；图片恢复按 exact generation_attempt、started audit、独立人工决策隔离不确定结果；Owner 角色变更在 advisory lock 下重新核对 actor；JSON Schema compiler 禁止外部 LoadURL；HTTP 错误脱敏、body 限制和 BFF TLS/签名 JWT 的测试是独立门禁。观察到 RecoverOwner 本地操作路径忽略 advisory-lock Exec 的返回值，但后续同事务查询仍检查错误；未把这一代码形态直接认定为已证实权限漏洞。此抽查与 115 个必要 DB 测试不等于对 780 文件完整安全认证。

[Issue #338](https://github.com/rushairer/gouno-blog/issues/338) 仍开放：action-by-action 恢复、跨存储与 Provider 结果限制仍须评估；早期 cancel/reclaim P0 已由 #356 修复，不应误报为当前未修漏洞。不确定执行不自动重放；没有声称 Provider exactly-once 或账单确定性。Connector 行为完全未修改，冻结约束及尚未审结的风险保留。

[PR #413](https://github.com/rushairer/gouno-blog/pull/413) 合并候选 manifest；[PR #414](https://github.com/rushairer/gouno-blog/pull/414) 通过完整门禁后合并为 `53d4a7c9eb781b12bc583800fdda3a62af515be2`，把 BFF gate 和生产 Gosso 示例固定到已独立验签的 v1.7.0 digest。首次 BFF 尝试在进入协议测试前，因 Alpine 已移除 tzdata=2026e-r0 而失败；根据官方 [x86_64](https://pkgs.alpinelinux.org/package/v3.24/main/x86_64/tzdata) / [aarch64](https://pkgs.alpinelinux.org/package/v3.24/main/aarch64/tzdata) 包记录，builder/runtime 精确升级为 2026f-r0，并同步精确 pin 合同测试。保留基础镜像 digest、musl floor、扫描阈值与 Connector 源码。重跑的两架构构建/扫描、完整 CI 和 BFF 全成功；失败的首次尝试没有被算作通过。

新的真实 TLS BFF 验收覆盖 Seed confidential client/resource/back-channel 策略、PKCE 与 Blog-only audience、UserInfo audience 拒绝、refresh audience 保持、HttpOnly flow handle/AAL2 请求、HTTPS Back-Channel Logout 和服务端 refresh credential revoke。旧 v1.6.1 结果没有被直接外推为 v1.7.0 通过。

Blog `blog-admin-ai` 仍是 `needs-manual-recertification`；其他认证的 verified 不覆盖该缺口。新的 v1.7.0 BFF 成功只认证本次候选提交，不代表旧 v1.5.26 Tag 已通过该新组合。Oxlint warning budget 为 61，不是零 warning。没有因这些绿色流水线创建 Blog Tag。

## Docker、签名、SBOM 与 Helm

Gosso v1.7.0 [Release Workflow](https://github.com/rushairer/gosso/actions/runs/42013419089) 的所有步骤 completed/success，包括签名后的精确身份验签和 Helm 发布。另使用本地官方 Cosign 3.1.4 完成独立验证。

| 产物 | 不可变 digest |
| --- | --- |
| ghcr.io/rushairer/gosso:v1.7.0 multi-arch index | `sha256:b88b0475af7633afa8d6fda8ca7c3f4297d235a36e88219d9010818e76ff889a` |
| linux/amd64 child manifest | `sha256:71b33234b42340092d4ed39f726854c89554e136cea2b9a7d463421744b62637` |
| linux/arm64 child manifest | `sha256:b0d232aa5a0ee257951f16da844ffbc1cd3fdb306340f7d69be3514f481a1079` |
| oci://ghcr.io/rushairer/gosso/charts/gosso:1.7.0 | `sha256:97a538f2049c377218051bcaa58b225c77495e1140a4193c646afcac6601f33c` |

`cosign verify` 与 `cosign verify-attestation --type cyclonedx` 均退出 0，绑定完整 identity `https://github.com/rushairer/gosso/.github/workflows/release.yml@refs/tags/v1.7.0` 与 issuer `https://token.actions.githubusercontent.com`，没有 regex identity 或跳过日志验证选项。签名 claims 与 attestation subject 均指向上述 index digest；本地验证确认 claims、透明日志包含证据及受信任 CA。透明日志包含证据由 Cosign 离线验证，不表述为重新提交签名或在线重新签发。

下载并解析两项正式附件，其实际 SHA256 与 GitHub asset digest 一致：

- [CycloneDX](https://github.com/rushairer/gosso/releases/download/v1.7.0/gosso-sbom.cdx.json)：`bc9f1bff0d16ad5b2a8e37912dc3004948354cb8bd486f79b2bed3c8a0c7dbae`，内容与已验证的 CycloneDX attestation predicate 完全一致。
- [SPDX 2.3](https://github.com/rushairer/gosso/releases/download/v1.7.0/gosso-sbom.spdx.json)：`8f73efe159ea68a9b2aeba70671d75a76ccc81faaa4bcec8a5007c51043a3175`。

另外从 OCI index 读取两种架构各自的 Buildx SPDX 2.3 文档，各含 102 个 packages；不将单份 CycloneDX 附件冒称为两套架构独立扫描。Helm 4.2.3 实际回读 published chart，name=gosso、version=appVersion=1.7.0。

保留可复核文件：[artifact verification receipt](gosso-v1.7.0-artifact-verification.json)、[Cosign claims](gosso-v1.7.0-cosign.json)、[OCI index manifest](gosso-v1.7.0-manifest.json)。

可直接重现正式产物验签与 Helm 回读：

```sh
GOSSO_AUDIT_IMAGE=ghcr.io/rushairer/gosso@sha256:b88b0475af7633afa8d6fda8ca7c3f4297d235a36e88219d9010818e76ff889a
GOSSO_AUDIT_IDENTITY=https://github.com/rushairer/gosso/.github/workflows/release.yml@refs/tags/v1.7.0
cosign verify --certificate-identity "$GOSSO_AUDIT_IDENTITY" --certificate-oidc-issuer https://token.actions.githubusercontent.com "$GOSSO_AUDIT_IMAGE"
cosign verify-attestation --type cyclonedx --certificate-identity "$GOSSO_AUDIT_IDENTITY" --certificate-oidc-issuer https://token.actions.githubusercontent.com "$GOSSO_AUDIT_IMAGE"
helm show chart oci://ghcr.io/rushairer/gosso/charts/gosso --version 1.7.0
```


另对正式 index digest 使用 Trivy 0.76.1 分别按 linux/amd64、linux/arm64 扫描 Alpine OS 与 Go binary，HIGH/CRITICAL 均为 0，两个命令均退出 0，未使用 ignore-unfixed。数据库 UpdatedAt=2026-11-05T07:00:54Z，DownloadedAt=2026-11-05T11:16:53Z。该按正式 digest 的新鲜双架构结果与 CI 本地构建扫描分开保存于 receipt。

重建限制：Gosso Dockerfile 仍使用可变的基础镜像 tag、apk upgrade 与构建时间注入；Blog 的 APK 精确 pin 仍依赖上游仓库保留。当前产物通过不可变 digest、签名与 SBOM 可再次验证，不代表任意未来时间重新构建都得到逐位相同 digest。后续应单独治理基础镜像与 APK 快照保留；本次没有重写已发布 Tag 来隐藏这个限制。

Core、Template、CLI、Doc、UI、Demo 不以 Docker 镜像作为本次发布产物。Admin v0.7.25 和 Blog v1.6.0 未打 Tag，其候选 release digest / release identity / SBOM 验收全部待完成，不使用 main-channel digest 冒充正式产物。Blog 已有 post-signing cosign verify 实现与负向行为测试；它仍须在实际正式 Tag 发布后独立验证三张镜像。

## 未完成事项与可复现续作

1. 续作必须先 fetch，重新核对远端 main、Tag 和 exact-SHA 所有必要 workflow；不能复用已失效的 CI 或取消结果。
2. Admin 的 v1.7.0 全局治理镜像验收已通过；在明确 staging 环境、manual-first 全状态 rendered review 与敏感路径 PR 完成前不发布 v0.7.25。不能用脚本自动把认证改成 verified。
3. Blog 等待 Admin 上游，再完成 `blog-admin-ai` 人工优先认证、全 diff 安全/迁移审核、最终 SHA 的完整门禁与 staging/回滚验收。更新候选 manifest 中所有 PENDING 后才能宣布正式产物验收完成。
4. Blog 原克隆历史 Tag 冲突保留待历史核查，禁止 force fetch 覆盖。正式 Tag 操作使用可独立确认的远端对象；不得借本次整理重写历史。
5. 没有部署任何生产主机；本任务没有明确生产目标与操作资料。生产部署需单独授权并读取 LOCAL-OPERATIONS.md。

重现已发布产物检查：

```sh
git ls-remote --tags https://github.com/rushairer/gouno refs/tags/v1.4.0 'refs/tags/v1.4.0^{}'
git ls-remote --tags https://github.com/rushairer/gouno-template refs/tags/v1.4.0 'refs/tags/v1.4.0^{}'
git ls-remote --tags https://github.com/rushairer/gouno-cli refs/tags/v1.3.0 'refs/tags/v1.3.0^{}'
GOWORK=off GOTOOLCHAIN=go1.27.3 GOPROXY=https://proxy.golang.org go mod download -json github.com/rushairer/gouno@v1.4.0 github.com/rushairer/gouno-cli@v1.3.0
npm view @gouno/ui@0.4.9 version dist.integrity dist.tarball --json
npm view @gosso/client@0.9.3 version dist.integrity dist.tarball --json
```

以下命令可重现**已通过的隔离兼容性验收**；当前脚本/工作流已固定到正式 Gosso v1.7.0 digest。最终候选 SHA 有相关变化时应在干净的 GitHub runner 重跑；不要在现存部署目录直接运行带 Compose down -v 清理的兼容脚本：

```sh
gh workflow run gosso-release-compat.yml -R rushairer/gosso-admin --ref main
gh workflow run gosso-release-bff-compat.yml -R rushairer/gouno-blog --ref main
gh run list -R rushairer/gosso-admin --workflow gosso-release-compat.yml
gh run list -R rushairer/gouno-blog --workflow gosso-release-bff-compat.yml
```

版本控制副本位于 Gouno Core 仓库的 `docs/quality/`；同目录 [JSON 快照](ecosystem-release-evidence-2026-11.json) 保留远端引用、Release、官方 Go Proxy/SumDB、npm integrity 与 completed/success 流水线回执。Core main 的快照 SHA 是本报告提交前的审计基点；报告本身的提交及其 CI 通过 Git 历史核对。

本次执行遵循 blog-release 的 “Stop at the first failed gate” 与仓库 manual-first 约束；这些待办不被计入成功。发布链的完整目标尚未达成。
