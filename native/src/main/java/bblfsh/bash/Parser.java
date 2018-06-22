package bblfsh.bash;

import com.ansorgit.plugins.bash.lang.BashVersion;
import com.ansorgit.plugins.bash.lang.lexer.BashElementType;
import com.ansorgit.plugins.bash.lang.lexer.BashLexer;
import com.ansorgit.plugins.bash.lang.parser.BashParser;
import com.ansorgit.plugins.bash.lang.parser.BashParserDefinition;
import com.ansorgit.plugins.bash.lang.parser.BashPsiBuilder;
import com.ansorgit.plugins.bash.lang.parser.FileParsing;
import com.ansorgit.plugins.bash.lang.parser.BashParser;
import com.intellij.core.CoreApplicationEnvironment;
import com.intellij.core.CoreProjectEnvironment;
import com.intellij.lang.ASTNode;
import com.intellij.lang.ParserDefinition;
import com.intellij.lang.PsiBuilder;
import com.intellij.lang.impl.PsiBuilderFactoryImpl;
import com.intellij.mock.MockProject;
import com.intellij.openapi.Disposable;
import com.intellij.psi.tree.IElementType;
import com.intellij.openapi.extensions.ExtensionsArea;
import com.intellij.openapi.extensions.Extensions;
import com.intellij.lang.MetaLanguage;

public class Parser {
    public static ASTNode parse(final CharSequence code) {
        final MockProject project = project();
        final BashParser parser = new BashParser(project, BashVersion.Bash_v4);
        final ParserDefinition parserDefinition = new BashParserDefinition();

        final PsiBuilder builder = builder(parserDefinition, code);
        final IElementType root = parserDefinition.getFileNodeType();

        return parser.parse(root, builder);
    }

    private static MockProject project() {
        final ExtensionsArea rootArea = Extensions.getRootArea();
        CoreApplicationEnvironment.registerExtensionPoint(
                rootArea, MetaLanguage.EP_NAME, MetaLanguage.class);

        final CoreApplicationEnvironment appEnv = new CoreApplicationEnvironment(new Disposable() {
            @Override
            public void dispose() {
            }});

        final CoreProjectEnvironment environment = new CoreProjectEnvironment(new Disposable() {
            @Override
            public void dispose() {
            }
        }, appEnv);

        return environment.getProject();
    }

    private static PsiBuilder builder(final ParserDefinition parserDefinition, final CharSequence code) {
        return new PsiBuilderFactoryImpl().createBuilder(
                parserDefinition,
                new BashLexer(),
                code
        );
    }
}
