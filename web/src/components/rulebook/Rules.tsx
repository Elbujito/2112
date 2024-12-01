import { Section } from '../../shared/layout/Section';
import { VerticalFeatureRowWithSteps } from '../../shared/features/VerticalFeatureRowWithSteps';

const Rules = () => (
  <Section
    styles={`py-6 max-w-screen-lg mx-auto px-3 `}
    title="Principe of the game"
    description="Each player embodies a space nation. The 2 players start the game with 3 sovereign satellites each. The goal is to bring the opponent to lose all of their satellites either in LEO or in GEO.
    To achieve this, players take turns drawing cards from their deck (the previously shuffled card deck) and play certain cards that allow them to launch satellites, command centers or take actions that will help weaken the opponent or protect themselves."
  >
    <VerticalFeatureRowWithSteps 
      title="Preparing the Game"
      description="Determine which of the 2 players will start playing first. Each player shuffles their deck and draws the first 7 cards. If a player is not satisfied with their hand, they shuffle the 7 cards back into their library, draw 7 new ones, and put one of them back under their library."
    />
    <VerticalFeatureRowWithSteps
      title="Setup"
      steps={<span>Place the 2 ELINT satellite cards and the 3 IMINT satellite cards in LEO in front of you (face-up). Place the 4 SATCOM satellite cards in GEO in front of you (face-up). These cards represent your operation support capabilities.</span>}
      description="Place the game board between the players.Place the sovereign cards at each end of the game board. Each sovereign card represents their operational support capabilities."
      reverse
    />
    

    <VerticalFeatureRowWithSteps
      title="The steps of the duel"
      steps={
        <div>
        <ul className="space-y-4 list-decimal list-inside dark:text-black-100">
         <li>At the beginning of your turn, <b>untap</b> your cards.</li>
         <li>Draw a card from the top of your deck.</li>
         <li>If present, move the card on the Launcher slot to one of the LEO or GEO zones.</li>
            <li>
                Play cards up to the limit of: 
                <ol className="pl-5 mt-2 space-y-1 list-decimal list-inside">
                    <li>One capacity card (Satellite or C2) by placing it respectively on the Launcher or Command Center slot.</li>
                    <li>One instant action card.</li>
                </ol>
            </li>
           <li>And now its time for combat! Engage with your capacities and attack with them.</li>
        </ul>
        <span>Your opponent can choose to block one of your satellites with more than one of their satellites - or they can choose not to block it.</span>
        </div>
      }
      description=""
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="The battle phase"
      description="When attacking, you can target an opponent's sovereign satellite with your own satellites. However, you cannot attack your opponent's satellite capabilities (cards) directly. Your opponent has the choice to block your attack with their own satellites or not."
      steps={
        <div>
        <ul className="space-y-4 list-none list-inside dark:text-black-100">
            <li>
              During your combat phase, you can choose to attack with one or multiple satellites
                <ol className="pl-5 mt-2 space-y-1 list-disc list-inside">
                    <li>To attack, a satellite must be untapped.</li>
                    <li>At the beginning of the game, only one satellite can attack at a time.</li>
                    <li>You must engage a Command Center for each additional attacking satellite.</li>
                </ol>
            </li>
        </ul>
        <span>During combat, each attacking or blocking satellite deals damage equal to its power. If a satellite receives damage equal to or greater than its resistance, it is destroyed. If the damage dealt is less than the resistance, the target satellite is not destroyed.</span>
        <span>If a satellite is unblocked, it deals its damage to the targeted sovereign satellite.</span>
        </div>
      }
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="C2"
      description="When you play a Command Center, it enters the battlefield tapped: you won't be able to use it until your next turn."
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="Stealthy - High Power Microwave"
      description="During an attack by a Stealthy satellite, and a block by a High-Powered Microwave, the High-Powered Microwave is unable to orient its effector to generate damage.
      During an attack by the High-Powered Microwave, and a block by a Stealth satellite, the Stealth satellite takes damage from the High-Powered Microwave and is destroyed."
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="High Power Microwave"
      description="During an attack with HPM, it reduces the resistance of each satellite participating in defense."
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="Obstruction"
      description="A satellite that has attacked during its turn cannot participate in blocking during the opponent's turn.
      Each attacking satellite can be blocked by a single defending satellite (defense does not require an available C2)."
      reverse
    />
    <VerticalFeatureRowWithSteps
      title="Counterattack"
      description="During an attack, the defending space nation applies at the end of the turn the sum of the powers engaged in defense to the satellites engaged in the attack."
      reverse
    />
  <VerticalFeatureRowWithSteps
      title="High Power Microwave - IMINT/ELINT"
      description="When using High Power Microwave  to attack in LEO, it only targets one of the satellites, IMINT or ELINT, as chosen by the attacker. The area effect does not apply to both satellites."
      reverse
    />



  </Section>
);

export { Rules };
