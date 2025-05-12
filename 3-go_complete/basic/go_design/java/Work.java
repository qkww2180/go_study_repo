public class Work {
    // 传一个interface进来，所以实际的行为是不确定的
    public static void read(Reader reader){
        reader.Read(8);
    }
    // 即使传一个具体的class进来，实际行为还是不确定的(因为class可以被继承和重写)
    public static void readFile(FileReader reader){
        reader.Read(8);
    }
    public static void main(String args[]){
        FileReader fr=new FileReader();
        read(fr);

        LogReader lr=new LogReader();
        readFile(lr);
     }
}

// javac -encoding UTF-8 .\basic\go_design\java\*.java && cd .\basic\go_design\java\ && java Work && cd -
