#### 1.编译目录指纹计算文件,并计算html目录的指纹集

> go build sha1.go;
>
> ./sha1-* html； 得到一个html.csv文件;

#### 2. 计算并比对指纹集

> 对html.csv文件 求Hash： `shasum -a 256 html.csv`;
>
> 比对上步骤得到的指纹和sha1-each-file-in-html-dir.\*.csv中的值是否相同；若不同，则代码有误；

#### 3.  确认公钥正确性

> 到公钥服务器获取公钥： gpg --keyserver hkps://keys.openpgp.org --search-keys B35A5FF1BB38B971FDD2F757882334389EBBE46B;
> 先确定内容与指纹(B35A5FF1BB38B971FDD2F757882334389EBBE46B)的一致性，然后比对获取内容和public-key.txt的一致性；若完全一致，则repo内的pubkey内容无误；

#### 4. 验证指纹文件的签名

> gpg --verify  sha1-each-file-in-html-dir.\*.csv.asc public-key.txt;
> 若签名无效，则判定源码无效；
