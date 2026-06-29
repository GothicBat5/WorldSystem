import javax.swing.*;
import java.awt.*;
import java.util.Random;

public class SourBeam 
{

    static Random random = new Random();
    static int windowsCreated = 0;
    static final int MAXQ = 15; // How many windows??

    public static void main(String[] args) 
    {
        

        SwingUtilities.invokeLater(() -> {

            createWindow("Main");

            Timer timer = new Timer(400, e -> {

                if (windowsCreated >= MAXQ) 
                {
                    ((Timer)e.getSource()).stop();
                    return;
                }

                createWindow("Window: " + (windowsCreated + 1));
            });

            timer.start();

        });

    }

    static void createWindow(String title) 
    {

        JFrame frame = new JFrame(title);
        frame.setSize(550, 560); // Window size !! 
        frame.setLocation(random.nextInt(900), random.nextInt(600));
        frame.setDefaultCloseOperation(JFrame.DISPOSE_ON_CLOSE);

        JPanel panel = new JPanel(new BorderLayout());
        JLabel label = new JLabel("* * * ERROR * * *", SwingConstants.CENTER);
        panel.add(label);
        frame.add(panel);
        frame.setVisible(true);
        windowsCreated++;
    }
}
