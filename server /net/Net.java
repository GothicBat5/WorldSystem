import java.net.URL;
import java.io.BufferedReader;
import java.io.InputStreamReader;

public class Net
{
    public static void main(String[] args) throws Exception
    {
        URL url = new URL("https://www.behance.net");

        BufferedReader reader = new BufferedReader(new InputStreamReader(url.openStream()));

        String line;

        while((line = reader.readLine()) != null)
        {
            System.out.println(line);
        }

        reader.close();
    }
}
