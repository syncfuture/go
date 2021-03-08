#git tag | foreach-object -process { git push origin --delete $_ }
#git tag | foreach-object -process { git tag -d $_ }
<<<<<<< HEAD
git tag v1.4.1
=======
git tag v1.4.0
>>>>>>> v1
git push
git push --tags