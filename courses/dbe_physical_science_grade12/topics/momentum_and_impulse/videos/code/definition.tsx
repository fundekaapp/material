import { SkiaCanvas } from "@remotion/skia";
import { useFont, Text, Fill } from "@shopify/react-native-skia";
import { staticFile, useVideoConfig } from "remotion";
import { AssetManager } from "./AssetManager";
import { z } from "zod";

const roboto = staticFile("Roboto-Bold.ttf");

export const HelloSkia: React.FC<z.infer<typeof helloSkiaSchema>> = ({
	color1,
	color2,
}) => {
	const { height, width } = useVideoConfig();

	const bigFont = useFont(roboto, 64);
	const smallFont = useFont(roboto, 30);

	if (bigFont === null || smallFont === null) {
		return null;
	}

	return (
		<SkiaCanvas height={height} width={width}>
			<AssetManager
				images={[]}
				typefaces={{
					Roboto: roboto,
				}}
			>
				<Fill color="black"/>

				<Text
					x={24}
					y={height / 2}
					font={bigFont}
					text="Momentum"
					color="#9584ff"
				/>

				<Text
					x={24}
					y={height / 2 + 64 + 8}
					font={smallFont}
					text="The linear momentum of a particle (object) is a vector quantity equal to the product of the mass of the particle (object) and its velocity."
					color="white"
				/>
			</AssetManager>
		</SkiaCanvas>
	);
};
