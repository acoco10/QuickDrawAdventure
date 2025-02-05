# Quick ReadyDraw Adventure
Quick ReadyDraw Adventure (QDA) combines a shoot 'em-up and a classic turn-based RPG.

I wanted to make a Western game that felt more like a goofy spaghetti western with one-liners, bombastic action, and some shonen vibes rather than something gritty or serious.

So far there is one building to explore and a battle that can be triggered by talking to the NPC.

Much of my time was spent developing the battle system as ebiten is fairly lightweight.

Gameplay:

QDA has two phases for battles: the dialogue and shooting phases.

The dialogue phase is meant to feel like the part in a Western where two gunslingers eye each other up, trade barbs, grimaces, and cigar smoke. Instead of directly attacking your opponent, you use skills like Insult or Stare Down to anger or intimidate to psyche them out and slow down their draw. It's a game of poker and when you click draw you're going all in and seeing who's faster. 

The shooting phase is like a classic style turn based-RPG with a few caveats. Whoever wins the draw gets to go first throughout the battle unless they use a skill that gives up their initiative. Instead of mana, you only have six shots before you need to use a turn to reload. Some skills use more or less ammunition so think strategically about when it might be best to reload! Integers are pretty small I haven't decided completely on how to handle leveling up but characters only have 7-10 health and do 1-6 damage per turn so things begin and end pretty quickly.

I hope to have the first "Chapter" finished soon where the story and gameplay can integrate in a cohesive way. But right now the game is mostly a prototype and tools for building an RPG in ebiten.
