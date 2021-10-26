# CHANGELOG

## Unreleased (2021-10-22)

### New feature

- **tradebot**: :beers: separate buy and sell end time, modify cond([`355f0dc`](https://gitlab.tocraw.com/root/toc_trader/commit/355f0dca120df7ffddf14f1e9ffc2c75f16bd23a)) (@TimHsu@M1BP-20210907)
- **debugee**: :twisted_rightwards_arrows: add debug configuration([`d779253`](https://gitlab.tocraw.com/root/toc_trader/commit/d7792539dbd4144f92a36f0b25cf2baaef58aff3)) (@TimHsu@M1BP-20210907)
- **exec**: :twisted_rightwards_arrows: add windows support, change to prompt, split buy sell wait time([`b7f1a0a`](https://gitlab.tocraw.com/root/toc_trader/commit/b7f1a0a277d8f6c7983ff9e168b9736ee4d3c185)) (@TimHsu@M1BP-20210907)
- **kbar**: :alien: add kbar, check target exist in new method, tmp simulate on windows([`551d6fb`](https://gitlab.tocraw.com/root/toc_trader/commit/551d6fbfa3ba4ae8d19337a5f4f9b07a6a774b14)) (@TimHsu@M1BP-20210907)
- **target**: :zap: change target find method, add don't discard over time trade([`7a09a3c`](https://gitlab.tocraw.com/root/toc_trader/commit/7a09a3c69a16c1e4bd47975b07663d2caf20ca07)) (@TimHsu@M1BP-20210907)
- **simulation**: :pencil2: add multi trade day simulate, modifiy fetch entire tick method([`18330ab`](https://gitlab.tocraw.com/root/toc_trader/commit/18330ab1279fda2cc1cb161b228607d6228d656f)) (@TimHsu@M1BP-20210907)
- **api**: :bento: modify trade cond switch api, add manual buy later api([`a195312`](https://gitlab.tocraw.com/root/toc_trader/commit/a195312c5bcd2f7b8309bfcb9f15ec01b8f69a1d)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **docker**: :pushpin: fix wrong config path again([`bcf9683`](https://gitlab.tocraw.com/root/toc_trader/commit/bcf9683f574aa3e58210a6a7e5ea0cbe437e900e)) (@TimHsu@M1BP-20210907)
- **docker**: :rocket: fix wrong config path([`f6222bd`](https://gitlab.tocraw.com/root/toc_trader/commit/f6222bd42eba2f7245ac12005c1377075b7b6da1)) (@TimHsu@M1BP-20210907)
- **target**: :white_check_mark: fix wrong target date([`0b06d83`](https://gitlab.tocraw.com/root/toc_trader/commit/0b06d83d0755807d41471ac51c9650fb80b2c0f0)) (@TimHsu@M1BP-20210907)
- **target**: :fire: fix no check day trade([`b86e111`](https://gitlab.tocraw.com/root/toc_trader/commit/b86e1111784b510ebcea052d81fef7a2a60191b6)) (@TimHsu@M1BP-20210907)
- **status**: :package: make sure first status back, calculate ticks per second([`8e5b64c`](https://gitlab.tocraw.com/root/toc_trader/commit/8e5b64c4cc043f48b11685d3a3a672394673830d)) (@TimHsu@M1BP-20210907)
- **stockclose**: :necktie: make sure update close success otherwise fullrestart([`d4c97bb`](https://gitlab.tocraw.com/root/toc_trader/commit/d4c97bbc799b0e99a1ec6854b0f0c9a2c23d84c7)) (@TimHsu@M1BP-20210907)
- **restart**: :chart_with_upwards_trend: change restart method test, change golang to 1.17.2([`cf0252e`](https://gitlab.tocraw.com/root/toc_trader/commit/cf0252e3a6e6fe60d6830292afc68830acf4e47d)) (@TimHsu@M1BP-20210907)

## v1.1.0 (2021-10-12)

### New feature

- **simulate**: :bookmark: disable switch by TSE001, parallel simulation([`a436f88`](https://gitlab.tocraw.com/root/toc_trader/commit/a436f88242a71bdd1b7b9856ba73af346fe45011)) (@TimHsu@M1BP-20210907)
- **simulate**: :white_check_mark: add estimate time, modify cond([`3fd4290`](https://gitlab.tocraw.com/root/toc_trader/commit/3fd429013d4ec9ab3b6eef64473b07c147a9ae56)) (@TimHsu@M1BP-20210907)
- **simulate**: :twisted_rightwards_arrows: save simulation data, improve simulate performance([`0d0dfab`](https://gitlab.tocraw.com/root/toc_trader/commit/0d0dfab1ab92a030a80076f86396dfad7cb8be3e)) (@TimHsu@M1BP-20210907)
- **tradebot**: :see_no_evil: add sell first([`e355335`](https://gitlab.tocraw.com/root/toc_trader/commit/e355335c6b523722099d5a6798c7ba1c1878a89a)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **tradebot**: :zap: fix short stock bug, re-simulate, add short stock switch by TSE([`2857d29`](https://gitlab.tocraw.com/root/toc_trader/commit/2857d29c0a8c4de7cefb2eaafb9ec45292302319)) (@TimHsu@M1BP-20210907)

## v1.0.0 (2021-10-03)

### New feature

- **simulation**: :arrow_down: add rsi gap to simulation([`b02e75d`](https://gitlab.tocraw.com/root/toc_trader/commit/b02e75da87d57b472d20c0ab65020315205670b6)) (@TimHsu@M1BP-20210907)
- **simulation**: :triangular_flag_on_post: add ticks period to simulation([`957bf81`](https://gitlab.tocraw.com/root/toc_trader/commit/957bf815751f4d1e0aa40427dfe3b6aad14aad98)) (@TimHsu@M1BP-20210907)
- **tradebot**: :zap: change quota method split fee discount, add toggle for simulate, change sell method([`f11a0a8`](https://gitlab.tocraw.com/root/toc_trader/commit/f11a0a8dfd025948910d0df414ec55b7953fbf3e)) (@TimHsu@M1BP-20210907)
- **ci**: :dizzy: artifacts expose as docker logs([`e7b430b`](https://gitlab.tocraw.com/root/toc_trader/commit/e7b430bccd95e74de357b907f9f644ceb36904ea)) (@TimHsu@M1BP-20210907)
- **ci**: :pencil2: add docker log as artifact([`98dc7d1`](https://gitlab.tocraw.com/root/toc_trader/commit/98dc7d1bc569f3d632fef87b648aab99b71894ad)) (@TimHsu@M1BP-20210907)
- **simulate**: :card_file_box: change method to simulate much more cond([`6417c54`](https://gitlab.tocraw.com/root/toc_trader/commit/6417c5412f73acc6a46c0ab57a57a346d314ff57)) (@TimHsu@M1BP-20210907)
- **simulate**: :boom: limit history count to 1300, restart time to 8:15, fix show TSE001([`40d1955`](https://gitlab.tocraw.com/root/toc_trader/commit/40d195570dec91f9ef3e8326cc76005ea438c46b)) (@TimHsu@M1BP-20210907)
- **main**: :construction_worker: add restar full service, modify all api fail message([`9b44704`](https://gitlab.tocraw.com/root/toc_trader/commit/9b44704f6d11ad7114237fcdfaa5260afc9ad3f3)) (@TimHsu@M1BP-20210907)
- **tradebot**: :pencil2: add tse001 snapshot([`f70bddb`](https://gitlab.tocraw.com/root/toc_trader/commit/f70bddb5a0b154d5ed6e60f77dd672b224a263f6)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **tradebot**: :poop: fix "Trun enable buy off" show even if it is off([`74d8e51`](https://gitlab.tocraw.com/root/toc_trader/commit/74d8e5159ad57438c0f47792fc1b533a1a27cdfa)) (@TimHsu@M1BP-20210907)
- **ci**: :bulb: remove expose as([`bbdaaa6`](https://gitlab.tocraw.com/root/toc_trader/commit/bbdaaa69914519f754fbde2847d8aaa59ee08d60)) (@TimHsu@M1BP-20210907)
- **ci**: :zap: try expose docker logs([`c48b7b0`](https://gitlab.tocraw.com/root/toc_trader/commit/c48b7b0577d889244d93116330f8b8139ca150b5)) (@TimHsu@M1BP-20210907)
- **ci**: :wrench: when docker color log is disable([`b1bf3a8`](https://gitlab.tocraw.com/root/toc_trader/commit/b1bf3a8d8d0810666404d44c869f75193d533463)) (@TimHsu@M1BP-20210907)
- **quota**: :wheelchair: fix wrong stock price([`b31ffbc`](https://gitlab.tocraw.com/root/toc_trader/commit/b31ffbcb01e4a6ecb7a7859ecebb3eab30ccea6f)) (@TimHsu@M1BP-20210907)
- **fullrestart**: :heavy_minus_sign: fix wrong cron, add default key([`5a893c2`](https://gitlab.tocraw.com/root/toc_trader/commit/5a893c21092702852c96a0deaa6d307026fafcb6)) (@TimHsu@M1BP-20210907)
- **tradebot**: :mute: fix TickAnalyzeCondition([`4057f24`](https://gitlab.tocraw.com/root/toc_trader/commit/4057f24e4538086e80ad5c58e0a09526d49dce69)) (@TimHsu@M1BP-20210907)
