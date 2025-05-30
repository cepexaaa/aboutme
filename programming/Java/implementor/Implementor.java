package info.kgeorgiy.ja.kubesh.implementor;

import info.kgeorgiy.java.advanced.implementor.Impler;
import info.kgeorgiy.java.advanced.implementor.ImplerException;
import info.kgeorgiy.java.advanced.implementor.tools.JarImpler;

import javax.tools.JavaCompiler;
import javax.tools.ToolProvider;
import java.io.BufferedWriter;
import java.io.File;
import java.io.IOException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.jar.Attributes;
import java.util.jar.JarOutputStream;
import java.util.jar.Manifest;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.jar.JarEntry;

/**
 * Implementation of {@link Impler} and {@link JarImpler} interfaces.
 * <p>
 * This class generates Java source code for implementations of specified interfaces.
 * The generated class has the same name as the input interface with the suffix "Impl".
 * Methods in the generated class return default values for their return types.
 * </p>
 *
 * @author Sergey Kubesh
 * @see Impler
 * @see JarImpler
 */
public class Implementor implements Impler, JarImpler {
    /**
     * Line separator for the current system.
     * <p>
     * This value is obtained using {@link System#lineSeparator()}.
     * </p>
     */
    private static final String LINE_SEPARATOR = System.lineSeparator();

    /**
     * Main method for running the Implementor from the command line.
     * <p>
     * Can accept:
     * <ul>
     *   <li>One argument: the fully qualified name of the interface to implement.</li>
     *   <li>Three arguments: flag {@code -jar}, the fully qualified name of the interface to implement, target jar-file</li>
     * </ul>
     * The generated implementation is saved in the directory corresponding to the package of the interface.
     * Also create jar-file by implemented interface if first argument is flagged {@code -jar}.
     *
     * @param args command line arguments. Expected: the fully qualified name of the interface or flag {@code -jar}, Interface, target jar-file
     * @throws IllegalArgumentException if the number of arguments is not equal to one or the argument is {@code null}.
     * @see #implement(Class, Path)
     */
    public static void main(String[] args) {
        if (!(args.length == 1 && args[0] != null) && !(args.length == 3 && args[0].equals("-jar"))) {
            System.err.println("Must input one interface");
            return;
        }
        try {
            if (args[0].equals("-jar")) {
                Class<?> clazz = Class.forName(args[1]);
                JarImpler jarImpler = new Implementor();
                if (args[2].toLowerCase().endsWith(".jar")) {
                    Path jarFile = Paths.get(args[2]);
                    jarImpler.implementJar(clazz, jarFile);
                } else {
                    System.err.println("Must input one interface or interface with jar-file");
                }
            } else {
                Class<?> clazz = Class.forName(args[0]);
                Path path = Path.of(clazz.getPackageName().replace('.', File.separatorChar));
                Impler impler = new Implementor();
                impler.implement(clazz, path);
            }
        } catch (ClassNotFoundException e) {
            System.err.println("Class not found: " + args[0]);
        } catch (ImplerException e) {
            System.err.println(e.getMessage());
        }
    }

    /**
     * Generates a Java class implementing the specified interface.
     * <p>
     * The generated class is saved in the specified directory. The class name is the same as the interface name
     * with the suffix "Impl". Methods in the generated class return default values for their return types.
     * </p>
     *
     * @param token the interface to implement.
     * @param root the root directory where the implementation should be saved.
     * @throws ImplerException if the implementation cannot be generated.
     * @throws NullPointerException if {@code token} or {@code root} is {@code null}.
     * @throws ImplerException if {@code token} is not an interface or is a private interface.
     * @see #generateImplementationInterface(Class)
     * @see #getFilePath(Class, Path)
     */
    @Override
    public void implement(Class<?> token, Path root) throws ImplerException {
        if (token == null || root == null) {
            throw new NullPointerException("Arguments must not be null");
        }
        if (!token.isInterface()) {
            throw new ImplerException("Only interfaces are supported");
        }

        if (Modifier.isPrivate(token.getModifiers())) {
            throw new ImplerException("Interface does not be private");
        }

        String implementationCode = generateImplementationInterface(token);
        Path filePath = getFilePath(token, root);

        try {
            Files.createDirectories(filePath.getParent());
            try (BufferedWriter writer = Files.newBufferedWriter(filePath)) {
                writer.write(implementationCode);
            }
        } catch (IOException e) {
            throw new ImplerException("Failed to write implementation to file: " + e.getMessage(), e);
        }
    }

    /**
     * Generates the source code for a class implementing the specified interface.
     * <p>
     * The generated class implements all non-private methods of the interface.
     * Methods return default values for their return types.
     * </p>
     *
     * @param implementor the interface to implement.
     * @return the generated Java source code as a {@link String}.
     * @see #realiseMethod(Method, StringBuilder)
     */
    private String generateImplementationInterface(Class<?> implementor) {
        StringBuilder codeOfInterface = new StringBuilder();
        String packageName = implementor.getPackage().getName();
        String simpleName = implementor.getSimpleName();
        String interfaceName = simpleName + "Impl";

        if (!packageName.isEmpty()) {
            codeOfInterface.append("package ").append(packageName).append(";").append(LINE_SEPARATOR).append(LINE_SEPARATOR);
        }

        codeOfInterface.append("public class ").append(interfaceName)
                .append(" implements ").append(implementor.getCanonicalName()).append(" {").append(LINE_SEPARATOR).append(LINE_SEPARATOR);

        for (Method method : implementor.getMethods()) {
            realiseMethod(method, codeOfInterface);
        }
        codeOfInterface.append("}").append(LINE_SEPARATOR);
        return codeOfInterface.toString();
    }

    /**
     * Generates the implementation of a single method from the interface.
     * <p>
     * The method is added to the generated class with a default return value.
     * </p>
     *
     * @param method the method to implement.
     * @param codeOfInterface the {@link StringBuilder} containing the generated class code.
     * @see #writeAllParameters(StringBuilder, Class[], boolean)
     */
    private void realiseMethod(Method method, StringBuilder codeOfInterface) {
        if (Modifier.isAbstract(method.getModifiers())) {
            codeOfInterface.append("\t@Override").append(LINE_SEPARATOR)
                    .append("\tpublic ").append(method.getReturnType().getCanonicalName())
                    .append(" ").append(method.getName()).append("(");

            Class<?>[] parameterTypes = method.getParameterTypes();
            Class<?>[] exceptionsOfMethod = method.getExceptionTypes();
            writeAllParameters(codeOfInterface, parameterTypes, true);
            codeOfInterface.append(") ");
            if (exceptionsOfMethod.length > 0) {
                codeOfInterface.append("throws ");
            }
            writeAllParameters(codeOfInterface, exceptionsOfMethod, false);
            codeOfInterface.append("{").append(LINE_SEPARATOR);
            if (method.getReturnType() != void.class) {
                codeOfInterface.append("\t\treturn ").append(getDefaultValue(method.getReturnType())).append(";").append(LINE_SEPARATOR);
            } else {
                codeOfInterface.append("\t\t// No return value").append(LINE_SEPARATOR);
            }
            codeOfInterface.append("\t}").append(LINE_SEPARATOR).append(LINE_SEPARATOR);
        }
    }

    /**
     * Writes parameters or exceptions to the generated code.
     * <p>
     * This method is used to write method parameters or exceptions to the generated class code.
     * </p>
     *
     * @param codeOfInterface the {@link StringBuilder} containing the generated class code.
     * @param parameterTypes the array of parameter or exception types.
     * @param isArguments if {@code true}, writes parameter names; else, writes only types.
     */
    private static void writeAllParameters(StringBuilder codeOfInterface, Class<?>[] parameterTypes, boolean isArguments) {
        for (int i = 0; i < parameterTypes.length; i++) {
            if (i > 0) {
                codeOfInterface.append(", ");
            }
            codeOfInterface.append(parameterTypes[i].getCanonicalName());
            if (isArguments) {
                codeOfInterface.append(" arg").append(i);
            }
        }
    }

    /**
     * Returns the file path for the generated implementation.
     * <p>
     * The file path is constructed based on the package of the interface and the root directory.
     * </p>
     *
     * @param token the interface to implement.
     * @param root the root directory where the implementation should be saved.
     * @return the {@link Path} to the generated file.
     */
    private Path getFilePath(Class<?> token, Path root) {
        String packageName = token.getPackage().getName();
        String relativePath = packageName.replace('.', File.separatorChar);
        String fileName = token.getSimpleName() + "Impl.java";

        return root.resolve(relativePath).resolve(fileName);
    }

    /**
     * Returns the default value for the specified type.
     * <p>
     * This method returns the default value for primitive types and {@code null} for another types.
     * </p>
     *
     * @param type the type for which to return the default value.
     * @return the default value as a {@link String}.
     */
    private String getDefaultValue(Class<?> type) {
        if (type == boolean.class) {
            return "false";
        } else if (type == byte.class || type == short.class || type == int.class || type == long.class) {
            return "0";
        } else if (type == float.class) {
            return "0.0f";
        } else if (type == double.class) {
            return "0.0";
        } else if (type == char.class) {
            return "'\\u0000'";
        } else {
            return "null";
        }
    }

    /**
     * Generates a JAR file containing the compiled implementation of the specified interface.
     * <p>
     * This method use {@code implement} to create implementation of the input interface.
     *   <ol>
     *     <li>Create local directory for temporary files.</li>
     *     <li>Compile implemented class.</li>
     *     <li>Create JAR file by them.</li>
     *     <li>Clear created directory.</li>
     *   </ol>
     *
     * @param token the interface to implement.
     * @param jarFile the target JAR file where the implementation will be stored.
     * @throws ImplerException if any error occurs during the implementation, compilation, or JAR creation process.
     * @throws NullPointerException if {@code token} or {@code jarFile} is {@code null}.
     * @see #implement(Class, Path)
     * @see #compile(List, List)
     * @see #createJar(Class, Path, Path)
     */
    @Override
    public void implementJar(Class<?> token, Path jarFile) throws ImplerException {
        if (token == null || jarFile == null) {
            throw new NullPointerException("Arguments must not be null");
        }

        Path tempDir = Paths.get("temp");
        try {
            Files.createDirectories(tempDir);
            implement(token, tempDir);

            Path generatedClass = getFilePath(token, tempDir);
            compile(List.of(generatedClass), List.of(token));

            createJar(token, tempDir, jarFile);
        } catch (IOException e) {
            throw new ImplerException("Failed to create JAR file: " + e.getMessage(), e);
        } finally {
            deleteDirectory(tempDir);
        }
    }

    /**
     * Make and fill Jar file by compiled implemented class.
     * <p>
     *   <ol>
     *     <li>Creates a manifest file with version {@code 1.0} and sets the main class to the implemented class.</li>
     *     <li>Adds the compiled class file to the JAR archive.</li>
     *   </ol>
     * </p>
     *
     * @param token the interface whose implementation is being archived.
     * @param tempDir the temporary directory containing the compiled class file.
     * @param jarFile the target JAR file where the implementation will be stored
     * @throws IOException if error occurs during the creation of the JAR file or while copying the class file.
     */
    private void createJar(Class<?> token, Path tempDir, Path jarFile) throws IOException {
        Manifest manifest = new Manifest();
        manifest.getMainAttributes().put(Attributes.Name.MANIFEST_VERSION, "1.0");
        manifest.getMainAttributes().put(Attributes.Name.MAIN_CLASS, token.getPackageName() + "." + token.getSimpleName() + "Impl");

        try (JarOutputStream jarOutputStream = new JarOutputStream(Files.newOutputStream(jarFile), manifest)) {
            Path classFile = tempDir.resolve(token.getPackageName().replace('.', File.separatorChar))
                    .resolve(token.getSimpleName() + "Impl.class");
            String entryName = token.getPackageName().replace('.', '/') + "/" + token.getSimpleName() + "Impl.class";
            jarOutputStream.putNextEntry(new JarEntry(entryName));
            Files.copy(classFile, jarOutputStream);
        }
    }

    /**
     * Delete directory
     * <p>
     *     This method walks through the directory and deletes all files and subdirectories.
     * </p>
     *
     * @param directory the directory to delete. If {@code null} or the directory does not exist, the method does nothing.
     */
    private void deleteDirectory(Path directory) {
        if (Files.exists(directory)) {
            try (Stream<Path> stream = Files.walk(directory)) {
                stream.sorted((a, b) -> -a.compareTo(b))
                        .forEach(path -> {
                            try {
                                Files.delete(path);
                            } catch (IOException e) {
                                System.err.println("Failed to delete file: " + path);
                            }
                        });
            } catch (IOException e) {
                System.err.println("Failed to delete directory: " + directory);
            }
        }
    }

    /**
     * Compile files.
     * <p>
     *     This method compiling files with dependencies in them.
     * </p>
     *
     * @param files the list of Java source
     * @param dependencies the list of classes representing the dependencies required for compilation.
     * @throws ImplerException if system java compiler is not available or if the compilation fails.
     */
    private static void compile(
            final List<Path> files,
            final List<Class<?>> dependencies
    ) throws ImplerException {
        final JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
        if (compiler == null) {
            throw new ImplerException("Could not find java compiler, include tools.jar to classpath");
        }
        final String classpath = getClassPath(dependencies).stream()
                .map(Path::toString)
                .collect(Collectors.joining(File.pathSeparator));
        final String[] args = Stream.concat(
                Stream.of("-cp", classpath, "-encoding", "UTF8"),
                files.stream().map(Path::toString)
        ).toArray(String[]::new);
        final int exitCode = compiler.run(null, null, null, args);
        if (exitCode != 0) {
            throw new ImplerException("Compiler exit code: " + exitCode);
        }
    }

    /**
     * Returns the classpath entries for the specified dependencies.
     * <p>
     * This method retrieves the locations of the JAR files or directories containing the specified classes.
     * These locations are used to set the classpath during compilation.
     * </p>
     *
     * @param dependencies the list of classes representing the dependencies.
     * @return a list of paths to the JAR files or directories containing the dependencies.
     * @throws RuntimeException if the location of a dependency cannot be determined.
     */
    private static List<Path> getClassPath(final List<Class<?>> dependencies) {
        return dependencies.stream()
                .map(dependency -> {
                    try {
                        return Path.of(dependency.getProtectionDomain().getCodeSource().getLocation().toURI());
                    } catch (final URISyntaxException e) {
                        throw new RuntimeException(e);
                    }
                })
                .toList();
    }
}
