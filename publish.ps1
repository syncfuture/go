#git tag | foreach-object -process { git push origin --delete $_ }
#git tag | foreach-object -process { git tag -d $_ }
git push
git tag v1.9.3
git push --tags