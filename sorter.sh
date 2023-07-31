#!/bin/bash

# sorter.sh 스크립트 내용 예시
# 이 스크립트는 입력으로 받은 CSV 파일을 정렬하는 것으로 가정합니다.
# 실제로는 원하는 작업을 수행하는 스크립트를 구현해야 합니다.

if [ $# -ne 1 ]; then
    echo "Usage: $0 <input_csv_file>"
    exit 1
fi

input_file=$1
output1="./export_tmp.csv"
#output_file="${input_file%.*}_sorted.csv"
output_file="sorted_result.csv"
rm -f $output_file $output1

# 여기에 실제로 CSV 파일을 처리하는 로직을 추가하세요.
# 이 예제에서는 간단히 파일을 정렬하여 새로운 파일로 저장하는 것으로 가정합니다.

#match=`echo $file |grep ".csv"`
#	if [[ $file == "" ]]; then
#		exit 0;
#	fi
#	if [[ $match == "" ]]; then
#		exit 0;
#	fi


## +된 항목 1차 가공
echo "" > $output1
fgrep -v "+" $input_file >> $output1
grep "+" $input_file | sed -e 's/,[a-zA-Z0-9]*+/,/g' -e 's/,[a-zA-Z0-9]*_D+/,/g' >> $output1
grep "+" $input_file | sed -e 's/+[a-zA-Z0-9]*,/,/g' -e 's/+[a-zA-Z0-9]*_D,/,/g' >> $output1

# print
cat $output1
echo "====================================================================="

## _D 항목 2차 가공
echo "" > $output_file
cat $output1 >> $output_file
grep "_D" $output1 | sed -e 's/_D//g' >> $output_file
sed -i  's/_D//g' $output_file
sed -i  's/8809086209954/CFSC0003/g' $output_file
sed -i  's/8809694621391/CFSC0049/g' $output_file
sed -i  's/8809694622831/CFSC0067/g' $output_file
sed -i  's/8809694622085/CFSC0085/g' $output_file
sed -i  's/8809086205604/CFSN0007/g' $output_file
sed -i  's/8809694620264/CFSC0004/g' $output_file
sed -i  's/ //g' $output_file

cat $output_file
#sort -t',' -k1 "$input_file" > "$output_file"

#echo "File sorted successfully. Output: $output_file"
#rm -f temp_file.csv  export.csv export_tmp.csv
