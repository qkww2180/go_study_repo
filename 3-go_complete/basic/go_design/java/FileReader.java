public class FileReader implements Reader{
    public byte[] Read(int n){
        System.out.println("this is FileReader");
        return new byte[n];
    }
}
