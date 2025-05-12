import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.InputStreamReader;
import java.nio.charset.Charset;

public class ReadFile{
    public static void read(String fileName) throws Exception{
        // try {
            File file = new File(fileName);   // throws  NullPointerException
            FileInputStream fis = new FileInputStream(file);  //throws  FileNotFoundException, SecurityException
            InputStreamReader isr = new InputStreamReader(fis, Charset.forName("UTF-8"));
            BufferedReader br = new BufferedReader(isr);
            String str;
            while ((str = br.readLine()) != null) {
                System.out.println(str);
            }
            System.out.println(str);
            br.close();
        // }catch (NullPointerException e) {
        //     System.out.println(e);
        // }catch (FileNotFoundException e) {
        //     System.out.println(e);
        // }catch(SecurityException e){
        //     System.out.println(e);
        // }catch (Exception e) {
        //     System.out.println(e);
        // }
    }
}