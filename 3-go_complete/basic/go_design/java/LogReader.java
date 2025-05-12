public class LogReader extends FileReader{
    public byte[] Read(int n){
        System.out.println("this is LogReader");
        byte[] array=super.Read(n);
        for(int i=0;i<n;i++){
            array[i]/=2;
        }
        return array;
    }
}
